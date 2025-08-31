#!/bin/bash

# fallback script, called in case the configured 'Startup Master Template' fails
# the script uses the CLI client to do basic configuration and make BlackJack reachable
# Additionally it serves as example on how to use the CLI client to configure BlackJack A.L.O.A.

# Enable USB functions RNDIS, CDC ECM
BlackJack_cli usb set --vid 0x1d6c --pid 0x1347 --manufacturer "MaMe82" --sn "deadbeef1337" --product "BlackJack by MaMe82" --rndis --cdc-ecm

# Configure USB ethernet interface "usbeth" to run a DHCP server
#   - use IPv4 172.16.0.1 for interface with netmask 255.255.255.252
#   - disable DHCP option 3 (router) by passing an empty value
#   - disable DHCP option 6 (DNS) by passing an empty value
#   - add a DHCP range from 172.16.0.2 to 172.16.0.2 (single IP) with a lease time of 1 minute
BlackJack_cli net set server -i usbeth -a 172.16.0.1 -m 255.255.255.248 -o "3:" -o "6:" -r "172.16.0.2|172.16.0.2|5m"

# Enable WiFi AP (reg US, channel 6, SSID/AP name: "BlackJack", pre shared key: "MaMe82-BlackJack", don't use nexmon firmware)
# Note: As a pre-shared key is given, BlackJack assume the AP should use WPA2-PSK
# Note 2: The SSID uses Unicode characters not necessarily supported by the console, but BlackJack supports UTF-8 ;-)
BlackJack_cli wifi set ap -r US -c 6 -s "💥🖥💥 Ⓟ➃ⓌⓃ🅟❶" -k "MaMe82-BlackJack" --nonexmon

# Configure USB ethernet interface "wlan0" to run a DHCP server
#   - use IPv4 172.24.0.1 for interface with netmask 255.255.255.0
#   - disable DHCP option 3 (router) by passing an empty value
#   - disable DHCP option 6 (DNS) by passing an empty value
#   - add a DHCP range from 172.24.0.10 to 172.24.0.20 with a lease time of 5 minutes
BlackJack_cli net set server -i wlan0 -a 172.24.0.1 -m 255.255.255.0 -o "3:" -o "6:" -r "172.24.0.10|172.24.0.20|5m"

BlackJack_cli led -b 2
