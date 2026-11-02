# Todo API Sözleşmesi

## Erişim adresleri

- API servisi: http://localhost:4000
- Gateway: http://localhost:8080

Aşağıdaki endpoint'lere her iki adres üzerinden erişilebilir.
Web UI, gateway üzerinden `/api/*` yollarını kullanır.

## Veri Modelleri

### 1. Todo Modeli

| Alan | Tip | Açıklama |
| --- | --- | --- |
| id | string | Sunucunun ürettiği benzersiz kimlik |
| title | string | Zorunlu; boş veya yalnızca boşluk olamaz |
| description | string | Opsiyonel; varsayılan boş metin |
| done | boolean | Varsayılan false |
| due_date | string | Opsiyonel; boş metin veya RFC3339 biçiminde tarih ve saat |
| visibility | string | `private` (varsayılan) veya `public` |
| user_id | string | Todo'yu oluşturan kullanıcının kimliği |
| username | string | Todo'yu oluşturan kullanıcının kullanıcı adı |
| created_at | string | Sunucunun belirlediği UTC oluşturulma zamanı |
| updated_at | string | Sunucunun belirlediği UTC son değişiklik zamanı |
| deleted_at | string veya null | UTC silinme zamanı; aktif kayıtta null |

### 2. User Modeli

| Alan | Tip | Açıklama |
| --- | --- | --- |
| id | string | Benzersiz kullanıcı kimliği |
| username | string | Benzersiz kullanıcı adı (en az 3 karakter) |
| role | string | `user` (standart kullanıcı) veya `superadmin` (süper admin) |
| created_at | string | UTC oluşturulma zamanı |

Varsayılan Sistem Hesapları:
- **admin** (Şifre: `admin123`, Rol: `superadmin`)
- **user1** (Şifre: `user123`, Rol: `user`)

## Kimlik Doğrulama ve Yetkilendirme Kuralları

API isteklerinde `Authorization: Bearer <token>` başlığı kullanılır.

### Görünürlük ve Erişim Kuralları:
1. **Private (`private`):** Yalnızca oluşturan kullanıcı ve **Süper Admin** tarafından görülebilir.
2. **Public (`public`):** Oturum açmış tüm kullanıcılar ve **Süper Admin** tarafından görülebilir.
3. **Güncelleme ve Silme:** Sadece todo'nun sahibi veya **Süper Admin** güncelleyebilir ve silebilir. Başka bir kullanıcının todo'sunu güncellemeye veya silmeye çalışıldığında `403 Forbidden` döner.
4. **Süper Admin (`superadmin`):** Tüm kullanıcıların hem gizli hem açık tüm todo'larını görüntüleyebilir, güncelleyebilir ve silebilir.

---

## Kimlik Doğrulama Endpoint'leri

### POST /api/auth/register

Yeni bir kullanıcı hesabı oluşturur ve oturum token'ı döner.

Request Body:
```json
{
  "username": "ahmet",
  "password": "guclu_sifre_123"
}
```

Response Body (201 Created):
```json
{
  "token": "a1b2c3d4...",
  "user": {
    "id": "u-123",
    "username": "ahmet",
    "role": "user",
    "created_at": "2026-11-02T12:00:00Z"
  }
}
```

### POST /api/auth/login

Kullanıcı adı ve şifre ile giriş yapar.

Request Body:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

Response Body (200 OK):
```json
{
  "token": "e5f6g7h8...",
  "user": {
    "id": "u-admin",
    "username": "admin",
    "role": "superadmin",
    "created_at": "2026-11-02T12:00:00Z"
  }
}
```

### GET /api/auth/me

Mevcut oturum açmış kullanıcının bilgilerini döner.

Headers:
`Authorization: Bearer <token>`

Response Body (200 OK):
```json
{
  "id": "u-admin",
  "username": "admin",
  "role": "superadmin",
  "created_at": "2026-11-02T12:00:00Z"
}
```

---

## Todo Endpoint'leri

### GET /api/todos

Kullanıcının görmeye yetkili olduğu aktif kayıtları listeler.
- Normal kullanıcı: Kendi todo'ları + diğer kullanıcıların `public` todo'ları.
- Süper Admin: Tüm kullanıcıların tüm todo'ları.

Headers:
`Authorization: Bearer <token>`

Response Body (200 OK):
```json
[
  {
    "id": "todo-123",
    "title": "Rapor hazırla",
    "description": "",
    "done": false,
    "due_date": "",
    "visibility": "private",
    "user_id": "u-123",
    "username": "user1",
    "created_at": "2026-11-02T12:00:00Z",
    "updated_at": "2026-11-02T12:00:00Z",
    "deleted_at": null
  }
]
```

### POST /api/todos

Yeni bir todo oluşturur.

Headers:
`Authorization: Bearer <token>`

Request Body:
```json
{
  "title": "Ekip toplantısı",
  "description": "Haftalık planlama",
  "visibility": "public"
}
```
- `visibility`: `"public"` veya `"private"` (belirtilmezse varsayılan `"private"`).

Response Body (201 Created):
```json
{
  "id": "todo-456",
  "title": "Ekip toplantısı",
  "description": "Haftalık planlama",
  "done": false,
  "due_date": "",
  "visibility": "public",
  "user_id": "u-123",
  "username": "user1",
  "created_at": "2026-11-02T12:10:00Z",
  "updated_at": "2026-11-02T12:10:00Z",
  "deleted_at": null
}
```

### PUT /api/todos/{id}

Var olan todo'yu günceller. Sadece sahibi veya Süper Admin güncelleyebilir.

Headers:
`Authorization: Bearer <token>`

Request Body:
```json
{
  "done": true,
  "visibility": "private"
}
```

Status Codes:
- `200 OK`: Başarılı güncelleme.
- `400 Bad Request`: Geçersiz JSON veya alan.
- `401 Unauthorized`: Giriş yapılmamış.
- `403 Forbidden`: Todo başkasına ait ve kullanıcı süper admin değil.
- `404 Not Found`: Todo bulunamadı veya silinmiş.

### DELETE /api/todos/{id}

Todo'yu soft delete ile siler. Sadece sahibi veya Süper Admin silebilir.

Headers:
`Authorization: Bearer <token>`

Status Codes:
- `204 No Content`: Başarılı silme.
- `401 Unauthorized`: Giriş yapılmamış.
- `403 Forbidden`: Todo başkasına ait ve kullanıcı süper admin değil.
- `404 Not Found`: Todo bulunamadı veya silinmiş.
