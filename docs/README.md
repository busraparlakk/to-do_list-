API servisi → port 4000
Web UI servisi → port 3000
Gateway → port 8080

## Projeyi Çalıştırma

Proje üç bağımsız Go servisinden oluşur:

| Servis | Klasör | Port | Görevi |
|---|---|---|---|
| API | api-service | 4000 | Todo işlemleri ve JSON dosyasında veri saklama |
| Web UI | web-ui-service | 3000 | HTML, CSS ve JavaScript dosyalarını sunma |
| Gateway | gateway | 8080 | İstekleri ilgili servise iletme |

Ön koşul: Go kurulu olmalıdır. Kontrol etmek için:

```powershell
go version
```

### Başlatma Sırası

Servisleri API → Web UI → Gateway sırasıyla, ayrı terminallerde başlatın.

Bu sıra, gateway'e istek geldiğinde arka servislerin hazır olmasını sağlar.
Gateway tek başına başlayabilir; ancak hedef servis kapalıysa ilgili isteğe
502 Bad Gateway yanıtı verir.

Aşağıdaki komutlar her terminalin proje kökünde
(`to-do_reverse_proxy`) açıldığını varsayar.

### 1. API Servisi

Birinci terminalde:

```powershell
cd api-service
go run .
```

API, http://localhost:4000 adresinde çalışır.
Veriler `api-service/data/todos.json` dosyasında saklanır.

Dosya yolu göreli olduğu için API'yi `api-service` klasöründen çalıştırın.

### 2. Web UI Servisi

İkinci terminalde:

```powershell
cd web-ui-service
go run .
```

Web UI, http://localhost:3000 adresinde çalışır.
`public` klasöründeki dosyaları sunar.

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

Üç terminali de açık bırakın. Bir servisi durdurmak için kendi
terminalinde Ctrl+C kullanın.

### Uygulamayı Açma

Tarayıcıda http://localhost:8080 adresini açın.

Web UI, API isteklerini `/api/todos` gibi göreli adreslerle gönderir.
Böylece tarayıcıdaki sayfa ve API istekleri aynı origin'i kullanır.
Tarayıcıdan doğrudan 3000 veya 4000 portuna gitmek gerekmez.

## Örnek İstekler

Aşağıdaki komutları dördüncü bir PowerShell terminalinde,
herhangi bir klasörden çalıştırabilirsiniz.

### API Sağlık Kontrolü

Bu istek, API'yi doğrudan kontrol etmek içindir:

```powershell
curl.exe -i http://localhost:4000/health
```

Beklenen: 200 OK ve `{"status":"ok"}`.

### Gateway Üzerinden Web UI

```powershell
curl.exe -i http://localhost:8080/
```

Beklenen: 200 OK ve HTML içeriği.

### Todo Listesi

```powershell
curl.exe -i http://localhost:8080/api/todos
```

Beklenen: 200 OK ve JSON listesi.

### Todo Oluşturma

PowerShell'deki tırnak aktarımı sorunlarını önlemek için JSON gövdesini
geçici bir dosyaya yazıp curl ile gönderin:

```powershell
$requestFile = Join-Path $env:TEMP "todo-create-request.json"
$body = @{ title = "Gateway testi" } | ConvertTo-Json
$utf8 = [System.Text.UTF8Encoding]::new($false)
[System.IO.File]::WriteAllText($requestFile, $body, $utf8)

curl.exe -i -X POST http://localhost:8080/api/todos -H "Content-Type: application/json" --data-binary "@$requestFile"

Remove-Item $requestFile
```

Beklenen: 201 Created ve oluşturulan todo.

### Tek Todo Getirme

POST yanıtındaki gerçek ID'yi aşağıya yazın:

```powershell
$todoId = "BURAYA_GERCEK_ID"
curl.exe -i "http://localhost:8080/api/todos/$todoId"
```

Beklenen: 200 OK ve ilgili todo.

### Todo Güncelleme

Aynı terminalde, gerçek ID'yi içeren `$todoId` değişkenini kullanın:

```powershell
$requestFile = Join-Path $env:TEMP "todo-update-request.json"
$body = @{ done = $true } | ConvertTo-Json
$utf8 = [System.Text.UTF8Encoding]::new($false)
[System.IO.File]::WriteAllText($requestFile, $body, $utf8)

curl.exe -i -X PUT "http://localhost:8080/api/todos/$todoId" -H "Content-Type: application/json" --data-binary "@$requestFile"

Remove-Item $requestFile
```

Beklenen: 200 OK, `done: true` ve yenilenmiş `updated_at`.

PUT endpoint'i yalnızca gönderilen alanları günceller.
Gönderilmeyen alanlar korunur.

### Todo Silme

```powershell
curl.exe -i -X DELETE "http://localhost:8080/api/todos/$todoId"
```

Beklenen: 204 No Content. Yanıt gövdesi boştur.

Silinen kaydı tekrar sorgulayın:

```powershell
curl.exe -i "http://localhost:8080/api/todos/$todoId"
```

Beklenen: 404 Not Found ve `Todo bulunamadı`.

## Backend Kapalıyken Davranış

1. API terminalinde Ctrl+C ile API'yi durdurun.
2. Web UI ve gateway çalışmaya devam etsin.
3. Aşağıdaki istekleri gönderin:

```powershell
curl.exe -i http://localhost:8080/api/todos
curl.exe -i http://localhost:8080/
```

Beklenen sonuçlar:

- `/api/todos` → 502 Bad Gateway
- `/` → 200 OK ve Web UI HTML'i
- Gateway terminali → API bağlantısının reddedildiğini belirten log

API'yi kendi klasöründen `go run .` ile yeniden başlatın.
Gateway'i yeniden başlatmadan `/api/todos` isteğini tekrarlayın.
Yanıt yeniden 200 OK olmalıdır.

## Gelecek Geliştirmeler — Henüz Uygulanmadı

- Health check: Backend servislerinin durumunu periyodik kontrol etme.
- Load balancing: İstekleri birden fazla backend instance'ına dağıtma.
- Config dosyası: Port ve hedef adreslerini kod dışından yapılandırma.
- Logging ve metrics: İstek yolu, hedef backend, süre ve sonuçları izleme.
- Middleware: Kimlik doğrulama, istek sınırlandırma ve header düzenleme.
- TLS/HTTPS: İstemci ile gateway arasındaki bağlantıyı şifreleme.
- Graceful shutdown: Kapanırken yeni istek kabulünü durdurup devam eden
  isteklere belirli bir süre içinde tamamlanma fırsatı verme.

## Geliştirme Yaklaşımı

Önce bağımsız test edilebilir parçalar oluşturulur. Her parça doğrulandıktan
sonra diğerleriyle birleştirilir. Başarılı işlemlerle birlikte hata
senaryoları da test edilir.