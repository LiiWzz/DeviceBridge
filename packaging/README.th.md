# ติดตั้ง DeviceBridge รุ่นทดลอง M01

รุ่นนี้เป็นโครงสร้างพื้นฐาน: Agent ยังไม่เชื่อมต่อ Relay และยังใช้ SSH ระยะไกลจริง
ไม่ได้ แต่สามารถติดตั้ง เปิดโปรแกรม ตรวจสถานะ และทดลองขั้นตอนดูแลเวอร์ชันได้

## ติดตั้งด้วยคำสั่งเดียวหลังแตกไฟล์

เลือกแพ็กเกจให้ตรงกับเครื่อง เช่น `linux-armv7`, `linux-arm64`, `linux-amd64`
หรือ CLI สำหรับ Mac `darwin-arm64` / `darwin-amd64` รับทั้งแพ็กเกจและ checksum
จากผู้เผยแพร่ที่เชื่อถือได้ งานนี้ยังไม่ได้เผยแพร่ GitHub Release สาธารณะ

ตรวจ checksum ก่อนแตกไฟล์: Linux ใช้ `sha256sum -c ชื่อไฟล์.tar.gz.sha256`
ส่วน Mac ใช้ `shasum -a 256 -c ชื่อไฟล์.tar.gz.sha256`
checksum ตรวจไฟล์เสียได้ แต่ไม่ยืนยันตัวตนผู้เผยแพร่หากผู้โจมตีเปลี่ยนได้ทั้งสองไฟล์

ตัวอย่าง CLI บน Mac ARM64 (เปลี่ยนชื่อให้ตรงกับแพ็กเกจที่ได้รับ):

```sh
tar -xzf devicebridge-cli-0.1.0-preview.4-darwin-arm64.tar.gz
cd devicebridge-cli-0.1.0-preview.4-darwin-arm64
sh install.sh
"$HOME/.local/devicebridge/cli/bin/devicebridge" help
```

ไม่ต้องใช้ `sudo` และไม่ต้องติดตั้ง Rust, Go, Docker, Python หรือ Node
บนอุปกรณ์ ARM ตัวติดตั้งใช้ shell และเครื่องมือไฟล์มาตรฐาน พร้อมเครื่องมือ SHA-256
หากเครื่องเป็น Buildroot/Yocto แบบปรับแต่ง ต้องตรวจว่าเปิด applet ที่จำเป็นไว้
รายละเอียดอยู่ใน [คู่มือฉบับเต็ม](README.md)

ตัวติดตั้งเลือกโฟลเดอร์ส่วนตัวตามส่วนประกอบโดยอัตโนมัติ:

| ส่วนประกอบ | คำสั่งหลังติดตั้ง |
|---|---|
| Agent | `"$HOME/.local/devicebridge/agent/bin/devicebridge-agent"` |
| Control API | `"$HOME/.local/devicebridge/control/bin/devicebridge-control"` |
| Relay | `"$HOME/.local/devicebridge/relay/bin/devicebridge-relay"` |
| CLI | `"$HOME/.local/devicebridge/cli/bin/devicebridge" help` |

Agent, Control และ Relay เปิดแบบ foreground; กด Ctrl-C เพื่อหยุด
Control ตรวจได้ที่ `http://127.0.0.1:8080/healthz` ส่วน Relay ยังไม่เปิด transport
การติดตั้งไม่เริ่ม service อัตโนมัติ ไม่แก้ PATH ไม่สร้างรหัสผ่าน และไม่ทับ config
เลือกโฟลเดอร์บน filesystem ที่อนุญาตให้รันโปรแกรมได้ หาก mount แบบ `noexec`
แม้คัดลอกไฟล์สำเร็จก็จะเปิดโปรแกรมไม่ได้ ตัวติดตั้งไม่รัน binary เพื่อทดลองแทนผู้ใช้
PostgreSQL สำหรับเครื่องพัฒนายังอยู่ใน [PR #15](https://github.com/LiiWzz/DeviceBridge/pull/15)
แยกจากแพ็กเกจ binary และไม่ต้องลงบนอุปกรณ์ ARM

## อัปเกรด ย้อนเวอร์ชัน ถอนติดตั้ง

อัปเกรดโดยแตกแพ็กเกจรุ่นใหม่แล้วรัน `sh install.sh` อีกครั้ง ใช้ `--prefix`
เดิมด้วยหากเคยระบุโฟลเดอร์เอง ตัวอย่างจัดการ Agent ที่ติดตั้งแล้ว:

```sh
sh "$HOME/.local/devicebridge/agent/manage.sh" --status
sh "$HOME/.local/devicebridge/agent/manage.sh" --rollback
# หยุด Agent ที่กำลังรันก่อนถอนติดตั้ง
sh "$HOME/.local/devicebridge/agent/manage.sh" --uninstall
```

การย้อนเวอร์ชันสลับกลับได้เฉพาะรุ่นปัจจุบันกับรุ่นก่อนหน้า โปรแกรมที่กำลังรัน
ต้องหยุดแล้วเปิดใหม่จึงจะใช้เวอร์ชันที่เปลี่ยน ตัวติดตั้งไม่ฆ่า process ให้
รุ่นเก่าถูกเก็บไว้จนถอนติดตั้ง จึงควรตรวจพื้นที่ flash ก่อนอัปเกรดหลายครั้ง
ไฟล์ config ของผู้ใช้ที่อยู่นอกโฟลเดอร์ `releases` จะถูกเก็บไว้เมื่อถอนติดตั้ง
หากพบไฟล์แปลกปลอมหรือ checksum ไม่ตรง ตัวจัดการจะหยุดเพื่อให้ตรวจสอบก่อน

การตรวจ OS/CPU ไม่ได้ยืนยัน libc หรือความเข้ากันได้กับทุก Linux distribution
โดยเฉพาะ Rust ที่ลิงก์ GNU libc ใช้แทน musl ไม่ได้โดยอัตโนมัติ ผลทดสอบบน Mac
หรือ Docker VM ไม่ใช่ผลบนอุปกรณ์ ARM RAM 128 MB จริง
