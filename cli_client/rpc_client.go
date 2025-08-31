package cli_client

import (
	"fmt"
	pb "github.com/theaayushpatel/BlackJack/proto"
	"errors"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)



func ClientConnectServer(rpcHost string, rpcPort string) (
	connection *grpc.ClientConn,
	client pb.BlackJackClient,
	ctx context.Context,
	cancel context.CancelFunc,
	err error) {
	// Set up a connection to the server.
	address := rpcHost + ":" + rpcPort
	//log.Printf("Connecting %s ...", address)
	connection, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to BlackJack RPC server: %v", err)
	}
	//defer conn.Close()

	client = pb.NewBlackJackClient(connection)

	// Contact the server
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	//defer cancel()

	err = nil
	return
}

func ClientTriggerGroupWait(host string, port string, groupname string, value int32) (err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { return errors.New(fmt.Sprintf("Could not connect to BlackJack RPC server: %v", err)) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx := context.Background()

	_,err = rpcClient.WaitTriggerGroupReceive(ctx, &pb.TriggerGroupReceive{GroupName: groupname, Value: value})

	return
}
func ClientTriggerGroupSend(host string, port string, groupname string, value int32) (err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { return errors.New(fmt.Sprintf("Could not connect to BlackJack RPC server: %v", err)) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx := context.Background()
	ctx,cancel := context.WithTimeout(ctx,TIMEOUT_SHORT)
	defer cancel()

	_,err = rpcClient.FireActionGroupSend(ctx, &pb.ActionGroupSend{GroupName: groupname, Value: value})

	return
}

func ClientCreateTempDir(host string, port string, dir string, prefix string) (resultPath string, err error) {
	return clientCreateTempDirOfFile(host,port,dir,prefix,true)
}

func ClientCreateTempFile(host string, port string, dir string, prefix string) (resultPath string, err error) {
	return clientCreateTempDirOfFile(host,port,dir,prefix,false)
}

func ClientMountUMSImage(host string, port string, file string, cdrom bool) (err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {return}
	defer connection.Close()
	client := pb.NewBlackJackClient(connection)
	_,err = client.MountUMSFile(
		context.Background(),
		&pb.GadgetSettingsUMS{
			File: file,
			Cdrom: cdrom,
		})
	return
}

func clientCreateTempDirOfFile(host string, port string, dir string, prefix string, dirOnlyNoFile bool) (resultPath string, err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {return}
	defer connection.Close()
	client := pb.NewBlackJackClient(connection)
	resp, err := client.FSCreateTempDirOrFile(
		context.Background(),
		&pb.TempDirOrFileRequest{
			Prefix: prefix,
			Dir: dir,
			OnlyFolder: dirOnlyNoFile,
		})
	if err != nil {return}
	resultPath = resp.ResultPath
	return
}


func ClientRegisterEvent(host string, port string,  evtType int64) (err error) {
	// open gRPC Client
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {return}
	defer connection.Close()
	client := pb.NewBlackJackClient(connection)
	evStream, err := client.EventListen(context.Background(), &pb.EventRequest{ListenType: evtType})
	if err != nil { return err }

	for {
		event, err := evStream.Recv()
		if err == io.EOF { break }
		if err != nil { return err }

		log.Printf("Event: %+v", event)
	}
	return nil
}

func ClientUploadFile(host string, port string, src io.Reader, folder pb.AccessibleFolder, filename string, forceOverwrite bool) (err error) {

	// open gRPC Client
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {return}
	defer connection.Close()
	client := pb.NewBlackJackClient(connection)

	//try to create remote file
	_, err = client.FSWriteFile(
		context.Background(),
		&pb.WriteFileRequest{
			Folder:folder,
			Filename: filename,
			Data: []byte{}, //empty chunk
			Append: false,
			MustNotExist: !forceOverwrite,
		})
	if err != nil {return}

	fmt.Printf("Start appending to '%s' in folder '%s'\n", filename, pb.AccessibleFolder_name[int32(folder)])

	// start appending chunks read from source file to remote file (Remote file is closed and opened every time, but
	// this avoids client to server streaming, which would be hard to implement for gRPC-web
	chunksize := 1024
	buf := make([]byte,chunksize)
	pos := int64(0)
	for {
		n,rErr := src.Read(buf)
		if rErr != nil {
			if rErr == io.EOF { break } else { return rErr }
		}

		sendData := buf[:n]
		client.FSWriteFile(
			context.Background(),
			&pb.WriteFileRequest{
				Folder:folder,
				Filename: filename,
				Data: sendData,
				Append: true,
				MustNotExist: false,
			})

		pos += int64(n)
	}

	return nil
}

func ClientGetLED(host string, port string) (ls *pb.LEDSettings, err error) {
	conn, client, ctx, cancel, err := ClientConnectServer(host, port)
	defer conn.Close()
	defer cancel()
	if err != nil { return }

	ls, err = client.GetLEDSettings(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("Error getting LED blink count: %v", err)
	}

	return
}

func ClientReboot(host string, port string, timeout time.Duration) (err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx := context.Background()
	if timeout > 0 {
		ctxNew,cancel := context.WithTimeout(ctx, timeout)
		ctx = ctxNew
		defer cancel()
	}

	_,err = rpcClient.Reboot(ctx, &pb.Empty{})
	return
}

func ClientShutdown(host string, port string, timeout time.Duration) (err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx := context.Background()
	if timeout > 0 {
		ctxNew,cancel := context.WithTimeout(ctx, timeout)
		ctx = ctxNew
		defer cancel()
	}

	_,err = rpcClient.Shutdown(ctx, &pb.Empty{})
	return
}


func ClientDeployGadgetSettings(host string, port string, newGs *pb.GadgetSettings) (gs *pb.GadgetSettings, err error) {
	conn, client, ctx, cancel, err := ClientConnectServer(host, port)
	defer conn.Close()
	defer cancel()
	if err != nil { return 	}


	gs, err = client.DeployGadgetSetting(ctx, newGs)
	if err != nil {
		log.Printf("Error deploying current USB Gadget Settings: %+v", err)
		//We have an error case, thus gs isn't submitted by the gRPC server (even if the value is provided)
		//in case of an error `gs`should reflect the Gadget Settings which are deployed, now that deployment of the
		//new settings failed. So we fetch the result manually
		gs, _ = client.GetDeployedGadgetSetting(ctx, &pb.Empty{}) //We ignore a new error this time, if it occures `gs` will be nil
	}

	return
}

func ClientGetDeployedGadgetSettings(host string, port string) (gs *pb.GadgetSettings, err error) {
	conn, client, ctx, cancel, err := ClientConnectServer(host, port)
	defer conn.Close()
	defer cancel()
	if err != nil {	return }

	gs, err = client.GetDeployedGadgetSetting(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("Error getting USB Gadget Settings count: %+v", err)
	}

	return
}

func ClientSetLED(host string, port string, ls pb.LEDSettings) (err error) {
	conn, client, ctx, cancel, err := ClientConnectServer(host, port)
	defer conn.Close()
	defer cancel()
	if err != nil { return }

	_, err = client.SetLEDSettings(ctx, &ls)
	if err != nil {
		log.Printf("Error setting LED blink count %d: %v", ls.BlinkCount, err)
	}


	return
}

func ClientDeployEthernetInterfaceSettings(host string, port string, settings *pb.EthernetInterfaceSettings) (err error) {
	// Set up a connection to the server.
	address := host + ":" + port
	//log.Printf("Connecting %s ...", address)
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to BlackJack RPC server: %v", err)
	}
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)

	// Contact the server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)
	defer cancel()

	_, err = rpcClient.DeployEthernetInterfaceSettings(ctx, settings)
	return

}

func ClientDeployWifiSettings(host string, port string, settings *pb.WiFiSettings) (state *pb.WiFiState, err error) {
	// Set up a connection to the server.
	address := host + ":" + port
	//log.Printf("Connecting %s ...", address)
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to BlackJack RPC server: %v", err)
	}
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)

	// Contact the server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)
	defer cancel()

	state,err = rpcClient.DeployWiFiSettings(ctx, settings)
	return
}

func ClientHIDRunScript(host string, port string, ctx context.Context, scriptPath string, timeoutSeconds uint32) (scriptRes *pb.HIDScriptResult, err error) {
	scriptReq := &pb.HIDScriptRequest{
		ScriptPath: scriptPath,
		TimeoutSeconds: timeoutSeconds,
	}

	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)

	scriptRes,err = rpcClient.HIDRunScript(ctx, scriptReq)
	return
}

func ClientHIDRunScriptJob(host string, port string, scriptPath string, timeoutSeconds uint32) (scriptJob *pb.HIDScriptJob, err error) {
	scriptReq := &pb.HIDScriptRequest{
		ScriptPath: scriptPath,
		TimeoutSeconds: timeoutSeconds,
	}

	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)
	defer cancel()

	scriptJob,err = rpcClient.HIDRunScriptJob(ctx, scriptReq)
	return
}

func ClientHIDCancelScriptJob(host string, port string, jobID uint32) (err error) {
	cancelReq := &pb.HIDScriptJob{
		Id: jobID,
	}

	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)
	defer cancel()

	_,err = rpcClient.HIDCancelScriptJob(ctx, cancelReq)
	return
}


func ClientHIDGetScriptJobResult(host string, port string, jobID uint32) (scriptRes *pb.HIDScriptResult, err error) {
	req := &pb.HIDScriptJob{
		Id: jobID,
	}

	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)
	//	defer cancel()

	scriptRes,err = rpcClient.HIDGetScriptJobResult(context.Background(), req)
	return
}

func ClientListTemplateType(timeout time.Duration, host string, port string, ttype pb.ActionDeploySettingsTemplate_TemplateType) (res []string, err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx := context.Background()
	if timeout > 0 {
		ctxNew,cancel := context.WithTimeout(ctx, timeout)
		ctx = ctxNew
		defer cancel()
	}

	switch ttype {
	case pb.ActionDeploySettingsTemplate_USB:
		ma,err := rpcClient.ListStoredUSBSettings(ctx, &pb.Empty{})
		if err != nil { return res,err }
		return ma.MsgArray,nil
	case pb.ActionDeploySettingsTemplate_TRIGGER_ACTIONS:
		ma,err := rpcClient.ListStoredTriggerActionSets(ctx, &pb.Empty{})
		if err != nil { return res,err }
		return ma.MsgArray,nil
	case pb.ActionDeploySettingsTemplate_WIFI:
		ma,err := rpcClient.ListStoredWifiSettings(ctx, &pb.Empty{})
		if err != nil { return res,err }
		return ma.MsgArray,nil
	case pb.ActionDeploySettingsTemplate_NETWORK:
		ma,err := rpcClient.ListStoredEthernetInterfaceSettings(ctx, &pb.Empty{})
		if err != nil { return res,err }
		return ma.MsgArray,nil
	case pb.ActionDeploySettingsTemplate_BLUETOOTH:
		ma,err := rpcClient.ListStoredBluetoothSettings(ctx, &pb.Empty{})
		if err != nil { return res,err }
		return ma.MsgArray,nil
	case pb.ActionDeploySettingsTemplate_FULL_SETTINGS:
		ma,err := rpcClient.ListStoredMasterTemplate(ctx, &pb.Empty{})
		if err != nil { return res,err }
		return ma.MsgArray,nil
	default:
		return res,errors.New("unknown template type")
	}


}

func ClientDeployTemplateType(timeout time.Duration, host string, port string, ttype pb.ActionDeploySettingsTemplate_TemplateType, name string) (err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx := context.Background()
	if timeout > 0 {
		ctxNew,cancel := context.WithTimeout(ctx, timeout)
		ctx = ctxNew
		defer cancel()
	}

	switch ttype {
	case pb.ActionDeploySettingsTemplate_USB:
		_,err = rpcClient.DeployStoredUSBSettings(ctx, &pb.StringMessage{Msg:name})
	case pb.ActionDeploySettingsTemplate_TRIGGER_ACTIONS:
		_,err = rpcClient.DeployStoredTriggerActionSetReplace(ctx, &pb.StringMessage{Msg:name})
	case pb.ActionDeploySettingsTemplate_WIFI:
		_,err = rpcClient.DeployStoredWifiSettings(ctx, &pb.StringMessage{Msg:name})
	case pb.ActionDeploySettingsTemplate_NETWORK:
		_,err = rpcClient.DeployStoredEthernetInterfaceSettings(ctx, &pb.StringMessage{Msg:name})
	case pb.ActionDeploySettingsTemplate_BLUETOOTH:
		_,err = rpcClient.DeployStoredBluetoothSettings(ctx, &pb.StringMessage{Msg:name})
	case pb.ActionDeploySettingsTemplate_FULL_SETTINGS:
		_,err = rpcClient.DeployStoredMasterTemplate(ctx, &pb.StringMessage{Msg:name})
	default:
		return errors.New("unknown template type")
	}

	return
}

func ClientDBBackup(timeout time.Duration, host string, port string, name string) (err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx := context.Background()
	if timeout > 0 {
		ctxNew,cancel := context.WithTimeout(ctx, timeout)
		ctx = ctxNew
		defer cancel()
	}

	_,err = rpcClient.DBBackup(ctx, &pb.StringMessage{Msg:name})

	return
}

func ClientDBRestore(timeout time.Duration, host string, port string, name string) (err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx := context.Background()
	if timeout > 0 {
		ctxNew,cancel := context.WithTimeout(ctx, timeout)
		ctx = ctxNew
		defer cancel()
	}

	_,err = rpcClient.DBRestore(ctx, &pb.StringMessage{Msg:name})

	return
}

func ClientDBList(timeout time.Duration, host string, port string) (names []string, err error) {
	address := host + ":" + port
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil { log.Fatalf("Could not connect to BlackJack RPC server: %v", err) }
	defer connection.Close()

	rpcClient := pb.NewBlackJackClient(connection)
	ctx := context.Background()
	if timeout > 0 {
		ctxNew,cancel := context.WithTimeout(ctx, timeout)
		ctx = ctxNew
		defer cancel()
	}

	backups,err := rpcClient.ListStoredDBBackups(ctx, &pb.Empty{})
	if err != nil { return names, err}

	return backups.MsgArray,nil
}
