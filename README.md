# **📡 Network Traffic Monitor**

A simple Go-based network traffic monitor for macOS that **alerts when network traffic is unusually high**.

---

## **🚀 Features**

✅ Monitors network traffic in real time  
✅ Logs updates **only when traffic changes**  
✅ Sends **macOS notifications** for high traffic  
✅ Can run in the background  
✅ **Automatically starts on boot (optional `launchd` daemon)**  

---

## **📦 Installation & Usage**

```sh
go mod tidy
go build -o netmon network_monitor.go
./netmon &
```

To stop the process:

```sh
ps aux | grep netmon
kill <PID>
```

---

## **⚙ Configuration**

Edit these constants in `network_monitor.go`:

```go
const (
    monitorInterval = 2 * time.Second // Check every 2 seconds
    thresholdMBps   = 1               // Alert if traffic > 1 MBps
    changeThreshold = 0.1             // Log only if traffic changes > 0.1 MBps
)
```

---

## **✅ 2. Create a macOS `launchd` Daemon**

This will **automatically start `netmon` at boot**, and you won’t need to start it manually.

### **🔹 Step 1: Create a `launchd` Service**

```sh
mkdir -p ~/Library/LaunchAgents
nano ~/Library/LaunchAgents/com.blyndon.networkmonitor.plist
```

Paste this inside:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>com.blyndon.networkmonitor</string>

    <key>ProgramArguments</key>
    <array>
      <string>/Users/YOUR_USERNAME/VSCodeProjects/network-monitor/netmon</string>
    </array>

    <key>RunAtLoad</key>
    <true/>

    <key>KeepAlive</key>
    <true/>

    <key>StandardOutPath</key>
    <string>/Users/YOUR_USERNAME/VSCodeProjects/network-monitor/network_monitor.log</string>

    <key>StandardErrorPath</key>
    <string>/Users/YOUR_USERNAME/VSCodeProjects/network-monitor/network_monitor.err</string>
  </dict>
</plist>
```

Replace **`YOUR_USERNAME`** with your actual macOS username.

### **🔹 Step 2: Load the Daemon**

```sh
launchctl load ~/Library/LaunchAgents/com.blyndon.networkmonitor.plist
launchctl start com.blyndon.networkmonitor
```

Now your **`netmon` runs automatically in the background**, even after a reboot! 🚀

### **🔹 Step 3: Check Status**

```sh
launchctl list | grep networkmonitor
```

### **🔹 Step 4: Stop or Remove It**

To **stop** it:

```sh
launchctl unload ~/Library/LaunchAgents/com.blyndon.networkmonitor.plist
```

To **remove it permanently**:

```sh
rm ~/Library/LaunchAgents/com.blyndon.networkmonitor.plist
```
