auth_service.go: Register bu fonksiyonda parolanın dolu ya da boş olduğunu kontrol etmemişsin. boş parola ile register olunabilir mi? 

bu fonksiyonda da boş parola için hash oluşturuluyor. boş parola kontrolü yok. 

```go
func (u *User) SetPassword(plainPassword string) {
	u.PasswordHash = hashPassword(plainPassword)
}
```

token ları go daki mutex ile yapman çok iyi. Bunlar genelde böyle yönetilmez ama bu da güzel bir yönetim şekli olmuş. session kontrolünde mutex i kitleyip kontrol etmen de güzel bir uyarı olmuş. 

Bir veri saklanırken iki farklı yerde saklanabilir: disk (db, json, yaml, txt) ve in-memory. 

Sen userları to-do ları disktte saklamışsın json olarak db görevi görmüş bundan yana sorun yok. Ama session u in-memory de saklaman biraz sorun yaratabilir. Uygulama kapanıp açıldığında yani bir şekilde sen docker ı kapattığında güncelleme yaptığında memory sıfırlanır ve oturum açmış kullanıcılar tekrar oturum açması gerekir. O yüzden session gibi değerleri de memory de tutmamak gerekebilir. Sesssion yönetimi çok ayrı bir konu sadece şimdilik in-memory / disk ayırımını bil yeter. Bu arada jwt kullansaydın in-memory problem yaratmazdı.


Hem login de hem register da auth_service.go da : 
```go
username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", User{}, errors.New("kullanıcı adı ve şifre zorunludur")
	}
```
Bu kod bloğu var. Burada da kod tekrarı yapmışsın. Bunun en basit yolu bir helper yazmak. user.go ya validate diye bir fonksiyon yazıp user.validate(req.body) diyip valide etmek. Ama bu en basit ve en amatör yolu. Bunun yerine sistemlerde artık dto vs. var. DTO kavramını araştır biraz. 

auth_middleware.go da güzel bir kod yazmışsın. Buradaki kod session yapısının nasıl saklandığını sana anlatmıştır inş. Çünkü python da, java da c# da ya da node.js de session.start diyorsun session un nasıl oluşturulduğunu nasıl kodlandığını anlamıyorsun burada diğer kodlama dillerinde görünmeyen fonksiyonu kendin yazmışın ve in-memory nin asıl yeri burası. middleware de in-memory odaklı bir yaklaşım sergilemişsin burada bir kaç değişiklik ile session u diske yazabilirdin. Genelde session yönetimi yapan kütüphaneler bunu input olarak alır. Örneğin: 
```python 
session.import('in-memory')
```
gibi. Bunlarla çok işimiz yok şimdilik ama neyin neden olduğunu anlamak için güzel bir proje olmuş. 


auth_service.go: 

```go
withAuth := func(h http.HandlerFunc) http.HandlerFunc {
		return authMiddleware(authService, h)
	}

    mux.HandleFunc("GET /api/auth/me", withAuth(authHandler.Me))
```

Buradaki kod bloğunun ne yaptığını / mantığını anlamak çok önemli. Bunu çoğu jr. yazılımcı bilmiyor. 

Burada bir eksik var gibi: 
```go

	mux.HandleFunc("GET /api/auth/me", withAuth(authHandler.Me))
	mux.HandleFunc("POST /api/todos", withAuth(createTodoHandler))
	mux.HandleFunc("GET /api/todos", withAuth(listTodosHandler))
	mux.HandleFunc("GET /api/todos/{id}", withAuth(getTodoHandler))
	mux.HandleFunc("PUT /api/todos/{id}", withAuth(updateTodoHandler))
	mux.HandleFunc("DELETE /api/todos/{id}", withAuth(deleteTodoHandler))

```

Burada withAuth kontrolü kullanıcı login olmuş mu onu kontrol ediyor. listTodosHandler fonksiyonuna baktığım kadarıyla sistemdeki tüm to-do ları getiriyor. Burada istenilen şuydu: admin tüm to-do ları görebilir diğer userlar tüm to-doları göremez sadece kendi to-do larını görebilir olması gerekiyordu. Login kontrolünü yapmışsın ama login olan kullanıcı kim ve hangi verilere erişebilir kontrolünü yapmamışsın. 

authn ve authz dediğimiz fark burada çıkıyor zaten:

authentication: kimlik doğrulama 
authorization: yetki doğrulama

Sen authentication yapmışsın kimlik doğrulamışsın ama yetki doğrulamamışsın. authn/authz farkını anlamak için çok iyi bir bug. Şuanki halinde authn iyi çalışıyor olabilir ama authz yok :D 

tüm handlelarında bu kod var: 
```go 
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
		return
	}
```
Burada iki problem var copy-paste yapacaksan neden tek bir yerde yazmıyorsun? tüm hadlelarda neden copy-paste kod var. Ve diğer en büyük problem: 
Kodları sırasıyla yazıyorum: 

```go 
main.go

// Kimlik doğrulaması gerektiren endpoint'ler
	withAuth := func(h http.HandlerFunc) http.HandlerFunc {
		return authMiddleware(authService, h)
	}

	mux.HandleFunc("GET /api/auth/me", withAuth(authHandler.Me))
	mux.HandleFunc("POST /api/todos", withAuth(createTodoHandler))
	mux.HandleFunc("GET /api/todos", withAuth(listTodosHandler))

auth_midleware.go

authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Yetkisiz erişim: Giriş yapmalısınız", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Geçersiz kimlik doğrulama başlığı biçimi", http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(parts[1])
		user, ok := authService.Authenticate(token)
		if !ok {
			http.Error(w, "Geçersiz veya süresi dolmuş oturum", http.StatusUnauthorized)
			return
		}



todo_handler.go

func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
		return
	}

```

auth_midle da kullanıcı login olmuş mu olmamış mı zaten kontrol ediyorsun. withAuth fonksiyonuyla bunu route a yani /api/todos a bağlıyorsun. Peki aynı kontrolü handler ın içinde neden yapıyorsun ki? kendi yazdığın koda mı güvenmiyorsun? aynı kontrolü birden fazla yapman kod kalitesini öldürür zaten ama ekstra performansı da bitirir. Ve yarın güncelleme yaptığında withAuth u kaldırırsın ama fonksiyon içinde de aynı kontrolü kaldırman gerekir vs. vs. 


Kodlarında mutlaka yorum satırı yaz. 