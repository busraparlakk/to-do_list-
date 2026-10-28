# Todo API Sözleşmesi

## Erişim adresleri

- API servisi: http://localhost:4000
- Gateway: http://localhost:8080

Aşağıdaki endpoint'lere her iki adres üzerinden erişilebilir.
Web UI, gateway üzerinden `/api/todos` yolunu kullanır.

## Veri modeli

| Alan | Tip | Açıklama |
| --- | --- | --- |
| id | string | Sunucunun ürettiği benzersiz kimlik |
| title | string | Zorunlu; boş veya yalnızca boşluk olamaz |
| description | string | Opsiyonel; varsayılan boş metin |
| done | boolean | Varsayılan false |
| due_date | string | Opsiyonel; boş metin veya RFC3339 biçiminde tarih ve saat |
| created_at | string | Sunucunun belirlediği UTC oluşturulma zamanı |
| updated_at | string | Sunucunun belirlediği UTC son değişiklik zamanı |
| deleted_at | string veya null | UTC silinme zamanı; aktif kayıtta null |

Tarih örneği: `2026-10-28T18:00:00Z`.

Eski dosya kayıtlarında `deleted_at` alanı bulunmuyorsa kayıt aktif
kabul edilir.

## İstek modelleri ve doğrulama

İstemciden gelen veriler DTO yapılarıyla alınır:

- CreateTodoRequest: title, description, due_date
- UpdateTodoRequest: title, description, done, due_date

id, created_at, updated_at ve deleted_at alanlarını sunucu yönetir.
Bu alanlar istek DTO'larının parçası değildir.

Oluşturma ve güncellemede ortak kurallar uygulanır:

- Başlığın başındaki ve sonundaki boşluklar temizlenir.
- Temizlendikten sonra boş kalan başlık reddedilir.
- due_date boş olabilir; doluysa geçerli bir RFC3339 tarihi olmalıdır.
- Geçersiz alan değerleri için 400 Bad Request döner.
- Doğrulama başarısız olursa değişiklik dosyaya yazılmaz.

JSON yanıtlarının Content-Type değeri:
`application/json; charset=utf-8`.

Hata yanıtları JSON değil, düz metindir.
Örnek hata mesajları:

```text
title boş olamaz
```

```text
due_date geçerli bir RFC3339 tarihi olmalı
```

## GET /api/todos

Yalnızca aktif kayıtları listeler. Silinmiş kayıtlar döndürülmez.

Request Body: Yok.

Response Body:

```json
[
  {
    "id": "todo-123",
    "title": "Reverse proxy öğren",
    "description": "Go ile reverse proxy geliştir",
    "done": false,
    "due_date": "2026-09-30T18:00:00Z",
    "created_at": "2026-09-25T15:00:00Z",
    "updated_at": "2026-09-25T15:00:00Z",
    "deleted_at": null
  }
]
```

Aktif kayıt yoksa `200 OK` ile boş JSON dizisi `[]` döner.

Status Codes:

- 200 OK
- 500 Internal Server Error

## GET /api/todos/{id}

Belirtilen kimliğe sahip aktif kaydı getirir.

Request Body: Yok.

Response Body:

```json
{
  "id": "todo-123",
  "title": "Reverse proxy öğren",
  "description": "Go ile reverse proxy geliştir",
  "done": false,
  "due_date": "2026-09-30T18:00:00Z",
  "created_at": "2026-09-25T15:00:00Z",
  "updated_at": "2026-09-25T15:00:00Z",
  "deleted_at": null
}
```

Kayıt yoksa veya silinmişse `404 Not Found` döner:

```text
Todo bulunamadı
```

Status Codes:

- 200 OK
- 404 Not Found
- 500 Internal Server Error

## POST /api/todos

Yeni bir todo oluşturur.

Request Body:

```json
{
  "title": "Reverse proxy öğren",
  "description": "Go ile reverse proxy geliştir",
  "due_date": "2026-09-30T18:00:00Z"
}
```

- title zorunludur.
- description ve due_date gönderilmezse boş metin olur.
- done sunucu tarafından false olarak belirlenir.
- deleted_at başlangıçta null olur.
- created_at ve updated_at başlangıçta aynı zamandır.
- Mevcut oluşturma işleminde DTO'da bulunmayan JSON alanları yok sayılır.

Response Body:

```json
{
  "id": "todo-123",
  "title": "Reverse proxy öğren",
  "description": "Go ile reverse proxy geliştir",
  "done": false,
  "due_date": "2026-09-30T18:00:00Z",
  "created_at": "2026-09-25T15:00:00Z",
  "updated_at": "2026-09-25T15:00:00Z",
  "deleted_at": null
}
```

Status Codes:

- 201 Created
- 400 Bad Request
- 500 Internal Server Error

## PUT /api/todos/{id}

Mevcut uygulamada kısmi güncelleme yapar:
yalnızca gönderilen alanlar değiştirilir, gönderilmeyen alanlar korunur.

Request Body:

```json
{
  "title": "Reverse proxy çalış",
  "description": "Go ile reverse proxy geliştir",
  "done": true,
  "due_date": "2026-09-30T18:00:00Z"
}
```

Yalnızca tamamlanma durumunu değiştirmek de mümkündür:

```json
{
  "done": true
}
```

Güncelleme kuralları:

- En az bir güncellenecek alan gönderilmelidir.
- Body tek bir JSON nesnesi içermelidir.
- DTO'da bulunmayan alanlar 400 Bad Request ile reddedilir.
- title gönderilmişse boş olamaz.
- description veya due_date, boş metin gönderilerek temizlenebilir.
- Mevcut DTO yapısında null değerler gönderilmemiş alan gibi değerlendirilir.
- done: false gönderilmesi geçerli bir güncellemedir.
- Başarılı güncellemede updated_at yenilenir.
- id ve created_at korunur.
- Geçersiz güncellemede kayıt ve updated_at değiştirilmez.
- Silinmiş kayıt güncellenemez.

Response Body:

```json
{
  "id": "todo-123",
  "title": "Reverse proxy çalış",
  "description": "Go ile reverse proxy geliştir",
  "done": true,
  "due_date": "2026-09-30T18:00:00Z",
  "created_at": "2026-09-25T15:00:00Z",
  "updated_at": "2026-09-25T16:00:00Z",
  "deleted_at": null
}
```

Geçerli bir body ile silinmiş veya olmayan kayda yapılan istek
`404 Not Found` döner:

```text
Todo bulunamadı
```

Status Codes:

- 200 OK
- 400 Bad Request
- 404 Not Found
- 500 Internal Server Error

## DELETE /api/todos/{id}

Aktif kaydı soft delete yöntemiyle silinmiş olarak işaretler.

Request Body: Yok.

Başarılı Response Body: Yok.

İşlem sırasında:

- Kayıt todos.json dosyasından kaldırılmaz.
- deleted_at alanına UTC silinme zamanı yazılır.
- updated_at yenilenir.
- Kaydın diğer alanları korunur.
- Başarılı işlem 204 No Content döner.

Silinen kayıt artık:

- GET /api/todos listesinde görünmez.
- GET /api/todos/{id} ile alınamaz.
- PUT /api/todos/{id} ile güncellenemez.

Aynı kayda tekrar DELETE gönderilirse 404 Not Found döner.
İlk silinme zamanı ve updated_at değiştirilmez.

Silinmiş veya olmayan kayıt için hata gövdesi:

```text
Todo bulunamadı
```

Status Codes:

- 204 No Content
- 404 Not Found
- 500 Internal Server Error

## Saklama ve geri yükleme kapsamı

Mevcut uygulamada silinmiş kayıtlar dosyada saklanır.

- Geri yükleme endpoint'i yoktur.
- Otomatik kalıcı silme işlemi yoktur.
- Saklama süresi belirlenmemiştir.

Saklama süresi, geri yükleme ve kalıcı silme kuralları ileride
ayrıca tasarlanacaktır.