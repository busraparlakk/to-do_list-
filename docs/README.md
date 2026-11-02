API servisi → port 4000
Web UI servisi → port 3000
Gateway → port 8080

## Projeyi Çalıştırma

Proje üç bağımsız Go servisinden oluşur:

| Servis | Klasör | Port | Görevi |
|---|---|---|---|
| API | api-service | 4000 | Kullanıcı yetkilendirme, Todo CRUD işlemleri ve JSON dosyasında veri saklama |
| Web UI | web-ui-service | 3000 | HTML, CSS ve JavaScript dosyalarını sunma |
| Gateway | gateway | 8080 | İstekleri ilgili servise iletme (Reverse Proxy) |

Ön koşul: Go kurulu olmalıdır. Kontrol etmek için:

```powershell
go version
```

### Başlatma Sırası

Servisleri API → Web UI → Gateway sırasıyla, ayrı terminallerde başlatın.

Aşağıdaki komutlar her terminalin proje kökünde (`to-do_reverse_proxy`) açıldığını varsayar:

### 1. API Servisi

Birinci terminalde:

```powershell
cd api-service
go run .
```

API, http://localhost:4000 adresinde çalışır.
Veriler `api-service/data/todos.json` ve `api-service/data/users.json` dosyalarında saklanır.

### 2. Web UI Servisi

İkinci terminalde:

```powershell
cd web-ui-service
go run .
```

Web UI, http://localhost:3000 adresinde çalışır. `public` klasöründeki dosyaları sunar.

### 3. Gateway

Üçüncü terminalde:

```powershell
cd gateway
go run .
```

Gateway, http://localhost:8080 adresinde çalışır.

Yönlendirme kuralları:
- `/api/` ile başlayan yollar → http://localhost:4000
- Diğer yollar → http://localhost:3000

---

## Uygulamayı Kullanma

Tarayıcıda **http://localhost:8080** adresini açın.

### Kullanıcı Hesapları ve Rol Yönetimi:
Uygulama ilk kez başlatıldığında varsayılan kullanıcılar otomatik oluşturulur:
- **admin** (Şifre: `admin123`): **Süper Admin** rolündedir. Tüm kullanıcıların oluşturduğu hem gizli (`private`) hem herkese açık (`public`) bütün todo'ları görüntüleyebilir, güncelleyebilir ve silebilir.
- **user1** (Şifre: `user123`): Standart kullanıcı rolündedir. Yalnızca kendi todo'larını ve diğer kullanıcıların herkese açık (`public`) paylaştığı todo'ları görebilir. Sadece kendi todo'larını güncelleyebilir veya silebilir.
- Arayüzden veya API üzerinden yeni kullanıcılar da kayıt olabilir (`Kayıt Ol` sekmesi).

### Gizlilik (Visibility) Özelliği:
- **Özel (Private):** Yalnızca todo'yu oluşturan kullanıcı ve süper admin görebilir.
- **Herkese Açık (Public):** Tüm oturum açmış kullanıcılar görebilir; ancak yalnızca oluşturan kullanıcı veya süper admin güncelleyebilir/silebilir.

---

## Örnek API İstekleri (cURL / PowerShell)

### 1. Kullanıcı Girişi ve Token Alma
```powershell
curl.exe -i -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d "{\"username\":\"user1\",\"password\":\"user123\"}"
```
Dönen yanıttaki `token` değerini sonraki isteklerde `Authorization: Bearer <TOKEN>` olarak kullanın.

### 2. Todo Oluşturma (Private veya Public)
```powershell
curl.exe -i -X POST http://localhost:8080/api/todos `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer <TOKEN>" `
  -d "{\"title\":\"Proje mimarisi çalış\",\"visibility\":\"private\"}"
```

### 3. Todoları Listeleme
```powershell
curl.exe -i http://localhost:8080/api/todos -H "Authorization: Bearer <TOKEN>"
```

### 4. Süper Admin ile Giriş
```powershell
curl.exe -i -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d "{\"username\":\"admin\",\"password\":\"admin123\"}"
```
Süper admin token'ı ile `GET /api/todos` çağrıldığında sistemdeki tüm kullanıcıların todo'ları listelenir.
