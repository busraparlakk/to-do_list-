                                       dışardan gelen istek   

                                                |
                                                v

                                           Go Reverse Proxy 
                                      path'e göre yönlendirme  

                        |                                                   |
                        v                                                   v

                    Web UI servisi                                         API servisi
                                                                        to-do crud backend
                                                                            
                                                                            |
                                                                            v

                                                                        todos.json

Go ile Reverse Proxy / API Gateway Projesi — Ayrıntılı Yol Haritası
Bu doküman sana kodu vermez; hangi dosyayı, hangi sırayla, neden o sırayla oluşturman gerektiğini anlatır. Her adımda üç şeyi bulacaksın:

Ne yapacaksın (somut, tek bir eylem)
Neden bu sırada / bu şekilde (mühendislik gerekçesi)
Nasıl doğrularsın (o adımı bitirdiğini nasıl anlarsın)
Mantık şu: her adım, bir öncekini test edilebilir kılmadan bir sonrakine geçmiyor. Bu, "önce her şeyi yaz, sonra hiçbir şey çalışmıyor, nerede hata var bilmiyorum" durumuna düşmemen için bilinçli bir sıralama.

Nihai Mimari (referans olarak elinde bulunsun)
Tarayıcı ──► :8080 (Go Reverse Proxy)
                 │
      path == /api/*  ──►  :4000 (API servisi)  ──►  todos.json
      path == diğer   ──►  :3000 (Web UI servisi)
Üç bağımsız process, üç bağımsız klasör, üç bağımsız çalıştırma komutu. Bunları önce ayrı ayrı sağlamlaştırıp en son birleştireceksin — bu proje boyunca tekrar edeceğin temel prensip bu. //bu 3 farklı dediği uı -apı- reverse proxy

FAZ 0 — İskelet ve Sözleşmeler (kod yazmadan önce)
Neden bu faz en başta?
Bir API'yi ve onu tüketen bir UI'yi aynı anda tasarlayarak kodlamaya başlarsan, ikisi arasında sürekli fikir değiştirirsin ("bu alanın adı due_date mi dueDate mi olacaktı?") ve bu, ilerlemeni yavaşlatan asıl şeydir. Bu yüzden önce sözleşmeyi (contract) kâğıt üstünde sabitliyoruz; kod bu sözleşmeye uyacak, sözleşmeyi kod belirlemeyecek.

Adım 0.1 — Proje klasör yapısını oluştur
todo-project/
  api-service/        (Faz 1'de dolduracaksın)
  web-ui-service/      (Faz 2'de dolduracaksın)
  gateway/             (Faz 3'te dolduracaksın)
  docs/
    CONTRACT.md
    README.md
Neden önce klasör yapısı? Üç servisin birbirinden fiziksel olarak ayrı klasörlerde yaşaması, sana şunu erken öğretir: bu üç şey birbirinden bağımsız deploy edilebilir, bağımsız çalıştırılabilir uygulamalardır. Aynı klasöre hepsini karıştırırsan, ileride "bu dosya hangi servise ait" karmaşası yaşarsın.

Doğrulama: tree komutuyla (ya da dosya gezgininde) yapıyı görebiliyor musun, kontrol et.

Adım 0.2 — docs/CONTRACT.md dosyasını oluştur
İçine şunları yaz (henüz kod değil, sadece tanım):

Veri modeli — bir todo nesnesinin alanları ve tipleri:
id: string
title: string (zorunlu, boş olamaz)
description: string (opsiyonel)
done: boolean (varsayılan false)
due_date: string (ISO 8601 formatında tarih, opsiyonel)
created_at: string, updated_at: string
Endpoint listesi, her biri için: HTTP metodu + path, beklenen request body, dönen response body, olası HTTP status kodları.
Neden bu dosya önce yazılır, koddan önce? Çünkü Faz 1'de API'yi yazarken de, Faz 2'de UI'yi yazarken de aynı dokümana bakacaksın. İki taraf da aynı kaynaktan besleniyor olacak; bu da ileride "API şunu döndürüyor sandım ama meğer öyle değilmiş" tipi hataları en aza indirir. Profesyonel ekiplerde buna API-first tasarım denir — implementasyondan önce sözleşmeyi dondurmak.

Doğrulama: Dokümanı bir arkadaşına (ya da kendine, bir gün sonra) verdiğinde, kod görmeden sadece bu dosyaya bakarak API'yi nasıl kullanacağını anlatabilmen gerekir.

burdayım Adım 0.3 — Port ve isimlendirme kararlarını docs/README.md'ye yaz
API servisi → port 4000
Web UI servisi → port 3000
Gateway (dışarıya açılan tek kapı) → port 8080
Neden bunu ayrı bir yere yazıyoruz? İleride (Faz 3'te) proxy'yi yazarken bu numaraları hafızandan değil, bu dosyadan kopyalayacaksın. Sabit değerleri (port numaraları gibi) bir yerde merkezi tutmak, sonradan değiştirmek istediğinde "kaç yerde bu sayı geçiyor" sorusuyla uğraşmanı engeller.

FAZ 1 — API Servisi (Backend)
Neden UI'den önce API?
Web UI, API olmadan test edilemez (veri gelecek bir yer yok) ama API, UI olmadan curl ile tamamen test edilebilir. Yani önce bağımsız test edilebilir olanı kurmak, seni daha az belirsizlikle ilerletir. Bu sıralama tesadüfi değil: "önce en az bağımlılığı olan katmanı sağlamlaştır" prensibi.
Adım 1.1 — Boş bir HTTP sunucusu ayağa kaldır (tek endpoint: /health)
Ne yapacaksın: api-service/ içinde, seçtiğin dilin en minimal HTTP sunucusunu kur. Sadece GET /health isteğine 200 OK ve {"status":"ok"} dönsün.

Neden CRUD'a değil de sağlık kontrolüne önce başlıyoruz? Çünkü şu an test etmek istediğin şey iş mantığı değil, "sunucu gerçekten ayakta mı, doğru portu dinliyor mu" sorusu. En küçük, en az mantık içeren parçadan başlayıp onu doğrulamak, ileride bir hata çıktığında "sorun temel kurulumda mı yoksa iş mantığında mı" ayrımını kolaylaştırır.

Doğrulama: curl http://localhost:4000/health çalıştır, 200 ve beklediğin JSON'u görmelisin.

Adım 1.2 — todos.json dosyasını ve okuma/yazma yardımcı fonksiyonlarını oluştur
Ne yapacaksın:

api-service/data/todos.json dosyasını oluştur, içi [] olsun.
Dosyadan okuma ve dosyaya yazma için iki ayrı yardımcı fonksiyon yaz (henüz hiçbir HTTP endpoint'ine bağlamadan — izole test edilebilir olmalı).
Neden bunu CRUD endpoint'lerinden önce, ayrı bir adım olarak yapıyoruz? Çünkü "dosyadan okuma/yazma" ve "HTTP isteğini işleme" iki ayrı sorumluluktur (buna yazılımda separation of concerns denir). Eğer bunları baştan karıştırırsan, bir hata aldığında "veri okunamadı mı yoksa HTTP route'u mu yanlış" ayrımını yapamazsın. Bu fonksiyonları HTTP'den bağımsız yazıp, mümkünse küçük bir script ile (endpoint açmadan) manuel çağırıp test et.

Doğrulama: Yardımcı fonksiyonu çağırıp elle bir todo nesnesi yazdır, dosyayı aç, JSON'un doğru göründüğünü gözle kontrol et.

Adım 1.3 — POST /api/todos endpoint'ini yaz
Neden CRUD'un ilk harfi olan Create ile başlıyoruz, Read ile değil? Çünkü listeleme (GET) endpoint'ini anlamlı test edebilmen için önce veri üretebiliyor olman lazım. Boş bir listeyi GET etmek seni hiçbir şey hakkında ikna etmez.

Ne yapacaksın:

Request body'den title, description, due_date oku.
title boşsa 400 Bad Request dön (validation) — bunu Faz 0'daki sözleşmene göre yap.
Yeni bir id üret (basit bir artan sayaç ya da UUID — hangisini seçtiğini ve nedenini not al).
created_at/updated_at alanlarını şu anki zaman ile doldur.
Adım 1.2'deki yazma fonksiyonunu çağırarak dosyaya ekle.
201 Created ve oluşturulan nesneyi dön.
Doğrulama: curl -X POST ile 3 farklı todo gönder, her seferinde farklı bir id döndüğünü gözle kontrol et, todos.json dosyasını açıp 3 kaydın orada olduğunu doğrula.

Adım 1.4 — GET /api/todos ve GET /api/todos/{id} endpoint'lerini yaz
Neden şimdi? Artık dosyada gerçek veri var (Adım 1.3 sayesinde), bu yüzden listeleme endpoint'i anlamlı bir şekilde test edilebilir.

Doğrulama:

GET /api/todos → Adım 1.3'te eklediğin 3 kaydı görmelisin.
GET /api/todos/{var olan id} → tek kaydı görmelisin.
GET /api/todos/{olmayan id} → 404 Not Found dönmeli. (Bu satırı atlama — hata durumlarını erken test etmek, mutlu senaryoyu test etmekten daha çok şey öğretir.)
Adım 1.5 — PUT /api/todos/{id} endpoint'ini yaz
Ne yapacaksın: Var olan kaydı bul, gönderilen alanlarla güncelle, updated_at'i yenile, 404 (olmayan id) ve 400 (geçersiz body) durumlarını ele al.

Doğrulama: Bir kaydı done: true yap, sonra GET ile tekrar çek, değişikliğin kalıcı olduğunu (dosyaya yazıldığını) doğrula.

Adım 1.6 — DELETE /api/todos/{id} endpoint'ini yaz
Doğrulama: Bir kaydı sil, GET /api/todos ile listede olmadığını, GET /api/todos/{id} ile 404 döndüğünü doğrula.
Adım 1.7 — CORS ayarını ekle
Neden CRUD bittikten sonra, en sona bırakıyoruz? Çünkü CORS, tarayıcının başka bir origin'den (Faz 2'de kuracağın Web UI, :3000'den :4000'e) gelen isteklere izin verip vermemesiyle ilgili bir tarayıcı güvenlik mekanizması. curl ile test ederken CORS hatası görmezsin (CORS sadece tarayıcıda uygulanır), bu yüzden bu noktaya kadar hiç sorun yaşamamış olabilirsin. Faz 2'de tarayıcıdan fetch çağırdığında aniden "CORS hatası" görürsen şaşırmaman için, bunu şimdiden bilerek ekliyoruz.

Doğrulama: Response header'larında Access-Control-Allow-Origin görünüyor mu, curl -I ile kontrol et.

Faz 1 Bitiş Kriteri (kendine sor)
 Tüm CRUD endpoint'leri curl ile çalışıyor mu?
 Hata senaryoları (404, 400) doğru status kodu dönüyor mu?
 todos.json dosyası her işlemden sonra beklediğim gibi güncelleniyor mu?
 Sunucuyu kapatıp tekrar açtığımda veri kayboluyor mu, kalıyor mu? (Kalması gerekir — dosya tabanlı DB'nin kazandırdığı şey bu.)
Bu kriterlerin hepsi karşılanmadan Faz 2'ye geçme — çünkü Faz 2'de bir hata gördüğünde, "acaba UI'de mi API'de mi" diye ayrım yapabilmen için API'nin sağlam olduğuna güvenmen gerekir.

FAZ 2 — Web UI Servisi
Neden statik dosya sunucusuyla başlıyoruz, form ile değil?
Aynı mantık: en az hareketli parçadan başla. Önce boş bir HTML sayfasının doğru porttan servis edildiğini doğrula, sonra üstüne davranış ekle.

Adım 2.1 — Statik dosya sunucusunu kur, boş bir index.html servis et
Ne yapacaksın: web-ui-service/ içinde :3000 portunu dinleyen, index.html, app.js, style.css dosyalarını servis eden minimal bir sunucu kur. index.html içine sadece "Merhaba" yaz, henüz to-do mantığı ekleme.

Doğrulama: Tarayıcıda http://localhost:3000 açtığında "Merhaba" yazısını görmelisin. Tarayıcı geliştirici konsolunda hiçbir hata olmamalı.

Adım 2.2 — Sabit (statik) HTML ile listeyi göster — henüz API'ye bağlanma
Ne yapacaksın: index.html içine elle, JavaScript'siz, 2-3 tane sahte (hardcoded) todo satırı yaz (başlık, checkbox, sil butonu görünsün ama hiçbiri çalışmasın).

Neden önce statik, sonra dinamik? Çünkü şu an test etmek istediğin şey görsel yapı ve CSS. Eğer bunu JavaScript/fetch mantığıyla aynı anda yaparsan, bir şey bozulduğunda "CSS mi bozuk, JS mi bozuk" ayrımını yapamazsın. Önce gözle görünen iskeleti sabitle.

Doğrulama: Sayfa düzeni (liste, checkbox, buton, form alanı) gözünle onayladığın gibi mi duruyor?

Adım 2.3 — app.js içinde GET /api/todos çağrısını yaz, listeyi API'den doldur
Ne yapacaksın: Sayfa yüklendiğinde fetch('http://localhost:4000/api/todos') çağır (şimdilik tam adres, proxy henüz yok — bunu neden yaptığımızı Faz 3'te değiştireceğiz), gelen veriyle Adım 2.2'deki statik HTML'i dinamik olarak oluştur.

Neden tam adres kullanıyoruz, göreli path değil? Çünkü şu an tarayıcı :3000'de, API :4000'de — göreli path (/api/todos) yazarsan tarayıcı bunu :3000/api/todos olarak yorumlar ve 404 alırsın. Bu, proxy'nin tam olarak neyi çözdüğünü ileride somut olarak görmen için bilinçli bir ara adım.

Doğrulama: Tarayıcı Network sekmesinde :4000/api/todos isteğinin gittiğini ve API'den dönen gerçek verinin ekranda göründüğünü doğrula. (CORS hatası alırsan Faz 1 Adım 1.7'ye geri dön.)

Adım 2.4 — Ekleme formunu POST /api/todos'a bağla
Doğrulama: Formdan bir todo ekle, sayfa yenilenmeden (ya da yenilenerek) listede görünmeli, todos.json dosyasında da gerçekten kayıtlı olmalı.

Adım 2.5 — Checkbox'ı PUT'a, sil butonunu DELETE'e bağla
Doğrulama: Uçtan uca: ekle → tamamlandı işaretle → sil. Her adımda hem ekranı hem todos.json dosyasını kontrol et, ikisi birbiriyle tutarlı olmalı.

Faz 2 Bitiş Kriteri
 Web UI, :4000'deki API'ye tam adresle bağlanıyor ve tüm CRUD işlemleri tarayıcıdan çalışıyor
 Sayfa yenilendiğinde veri kaybolmuyor (çünkü API'den her seferinde taze çekiliyor)
 Konsolda hata yok

FAZ 3 — Go Reverse Proxy / API Gateway (asıl hedef)
Buraya kadarki her şey, bu fazı anlamlı kılmak için yapıldı: elinde, birbirinden bağımsız, doğruluğu kanıtlanmış iki servis var. Şimdi onları tek kapıdan geçirecek üçüncü, bağımsız bir Go uygulaması yazacaksın.

Neden proxy'yi en son yazıyoruz?
Çünkü proxy, kendi başına test edilemeyen bir bileşendir — arkasında gerçek, çalışan servisler olmadan proxy'nin doğru çalışıp çalışmadığını anlayamazsın. Bu yüzden onu, bağımlılıkları (Faz 1 ve Faz 2) hazır olduktan sonra yazmak, "hata proxy'de mi yoksa arkadaki serviste mi" belirsizliğini ortadan kaldırır.
## burda kaldım

Adım 3.1 — Go modülünü başlat, boş bir "hello world" HTTP sunucusu yaz
Ne yapacaksın: gateway/ klasöründe go mod init çalıştır, :8080 portunu dinleyen ve her isteğe sabit bir metinle cevap veren minimal bir main.go yaz.

Neden proxy mantığından önce bu? Aynı prensip: önce "bu process doğru portta ayakta mı" sorusunu izole et. Go'da HTTP sunucusu kurmanın temel API'sini (http.HandleFunc, http.ListenAndServe) burada, hiçbir proxy karmaşası olmadan öğren/doğrula.

Doğrulama: curl http://localhost:8080 sabit metni dönüyor mu?

Adım 3.2 — Tek bir path için, tek bir backend'e yönlendiren en basit proxy'yi yaz
Ne yapacaksın: Sadece /api/ ile başlayan istekleri :4000'e yönlendiren bir handler yaz. Go'nun standart kütüphanesindeki net/http/httputil paketinin ReverseProxy tipini araştır ve kullan. Henüz Web UI'yi bağlama — tek bir yönlendirme kuralını sonuna kadar çalışır hale getirmeden ikinciye geçme.

Neden önce tek kural, sonra ikinci kural? Çünkü iki kuralı aynı anda yazıp bir hata aldığında, hatanın hangi kuralda olduğunu ayıramazsın. Bir kuralı tam anlamıyla doğrulayıp güvendikten sonra, ikinciyi eklemek, deseni tekrar etmek (aynı yapıyı kopyalamak) haline gelir — bu da hem daha az hataya hem de daha hızlı ilerlemeye yol açar.

Doğrulama: curl http://localhost:8080/api/todos → API servisinden gelen gerçek veriyi görmelisin (:4000'e hiç dokunmadan, sadece :8080 üzerinden).

Adım 3.3 — İkinci path kuralını ekle: geri kalan her şey Web UI'ye
Ne yapacaksın: /api/ dışındaki tüm path'leri :3000'e yönlendiren ikinci bir handler ekle. Go'nun http.ServeMux'ını (ya da tercih ettiğin router'ı) kullanarak iki kuralı bir arada yönet.

Neden bu sırayla (spesifik path önce, genel path sonra)? Çoğu router'da (Go'nun ServeMux'ı dahil) eşleşme kuralları en spesifikten en genele doğru düşünülmelidir; /api/ gibi belirli bir prefix'i, hiçbir şeyle eşleşmeyen "hepsi" (/) kuralından önce tanımlamazsan, genel kural her isteği yutabilir ve /api/ istekleri asla kendi handler'ına ulaşmaz. Bunu bilerek, önce dar kuralı, sonra geniş kuralı ekliyoruz.

Doğrulama:

curl http://localhost:8080/ → Web UI'nin index.html'i dönmeli.
curl http://localhost:8080/api/todos → hâlâ API'den veri dönmeli (bir önceki kuralı bozmadığını doğrula).
Adım 3.4 — Web UI'deki fetch çağrılarını tam adresten göreli path'e çevir
Ne yapacaksın: Faz 2'de http://localhost:4000/api/todos yazdığın her yeri /api/todos olarak değiştir (Adım 2.3'teki bilinçli ara adımı şimdi kapatıyorsun).

Neden bu değişikliği en son yapıyoruz, proxy tamamlandıktan sonra? Çünkü bu değişikliğin doğru çalıştığını ancak proxy zaten /api/'yi doğru yere yönlendiriyorsa görebilirsin. Sırayı tersine çevirip önce path'i değiştirseydin, proxy henüz hazır olmadığı için hangi hatanın neden kaynaklandığını (UI mı, proxy mı) karıştırırdın.

Doğrulama: Tarayıcıda sadece http://localhost:8080 adresine git (3000 ve 4000'i tarayıcıya hiç yazma). Ekleme, listeleme, güncelleme, silme işlemlerinin hepsinin uçtan uca, tek port üzerinden çalıştığını doğrula. Network sekmesinde tüm isteklerin :8080'e gittiğini gözle kontrol et.

Adım 3.5 — Hata senaryosunu gözlemle (bilinçli olarak bir backend'i kapat)
Ne yapacaksın: API servisini (:4000) durdur, Web UI'yi açık bırak, proxy'den /api/todos iste.

Neden bunu bilerek yapıyoruz? Gerçek dünyada backend'ler her zaman ayakta olmayabilir; proxy'nin bu durumda nasıl davrandığını (muhtemelen 502 Bad Gateway ya da bağlantı hatası) görmek, ileride Faz 5'teki "health check" ve "graceful error handling" konularına neden ihtiyaç duyacağını somut olarak anlamanı sağlar. Bu adımı atlarsan, bu ihtiyacı sadece teoride bilirsin, hissetmezsin.

Doğrulama: Aldığın hatayı not al (status kodu, mesaj). API'yi tekrar başlat, isteğin normale döndüğünü doğrula.

Faz 3 Bitiş Kriteri
 Tarayıcıda tek adres (:8080) yeterli, başka hiçbir port bilmene gerek yok
 /api/ ile başlayan her şey doğru şekilde API'ye gidiyor
 Diğer her şey doğru şekilde Web UI'ye gidiyor
 Bir backend kapandığında proxy'nin verdiği hatayı gözlemledim ve anladım
FAZ 4 — Doğrulama ve Belgeleme
Adım 4.1 — docs/README.md'yi güncelle
Şunları ekle: üç servisi hangi sırayla başlatman gerektiği (API → Web UI → Gateway — neden bu sıra? Çünkü proxy, arkasında zaten çalışan servisler bekliyor; onları önce ayağa kaldırmazsan proxy'nin verdiği hatalar yanıltıcı olur), her biri için çalıştırma komutu, örnek curl komutları.

Adım 4.2 — Sıfırdan tam senaryo testi
Bilgisayarını yeniden başlatmış gibi düşün: üç servisi sırasıyla başlat, tarayıcıda sadece :8080'i kullanarak bir kullanıcı gibi uçtan uca test et (ekle, tamamla, düzenle, sil).

Neden bunu ayrı bir adım olarak yapıyoruz? Geliştirme sırasında muhtemelen servisleri hiç kapatmadan, sürekli açık tutarak ilerledin. "Sıfırdan başlatma" senaryosu, gizli kalmış bağımlılıkları (örneğin "aslında A servisi açık değilken B servisi hata veriyor ama fark etmemiştim" gibi) ortaya çıkarır.

FAZ 5 — Sonraki Adımlar (şimdilik uygulama, sadece not al)
Faz 3 Adım 3.5'te gözlemlediğin hata senaryosu, bu listedeki ilk maddenin neden var olduğunu zaten sana gösterdi. Listenin geri kalanı da benzer şekilde, bugün karşılaştığın küçük sürtünmelerin büyütülmüş, üretim seviyesindeki karşılıkları:

Health check — proxy, backend'lerin ayakta olup olmadığını periyodik kontrol etsin (Adım 3.5'te elle gördüğün hatayı otomatik tespit etme).
Load balancing — aynı path için birden fazla backend instance'ı arasında dağıtım.
Config dosyası — Adım 0.3'te sabitlediğin port numaralarını artık kodda değil, bir config.yaml/.json dosyasında tutmak (böylece portu değiştirmek için kodu tekrar derlemek gerekmez).
Logging & metrics — her isteğin hangi path'e, hangi backend'e, ne kadar sürede yönlendirildiğinin kaydı.
Middleware katmanı — auth (API key/JWT), rate limiting, header manipülasyonu.
TLS/HTTPS — dışarıya açılan tek portu şifreli sunmak.
Graceful shutdown — proxy kapanırken devam eden istekleri tamamlayıp öyle kapanması.
Bu maddelerden birini seçip üzerine konuşmak istediğinde, aynı mantıkla (önce neden ihtiyaç var, sonra en küçük adımdan başlayarak) ayrı bir yol haritası çıkarabiliriz.

Bu Yol Haritasının Arkasındaki Genel Prensip
Fark etmiş olabilirsin: her fazda tekrar eden bir kalıp var —

En az bağımlılığı olan, en izole test edilebilir parçadan başla.
Bir parçayı tam doğrulamadan bir sonrakine geçme.
Bilerek "eksik" bir ara adımdan geçip (örn. tam adresle fetch), sonra o eksiği kapatarak (göreli path'e geçiş), neyin neyi çözdüğünü somut olarak gör.
Hata senaryolarını da mutlu senaryolar kadar erken test et.
Bu dört madde, bu projeye özgü değil — herhangi bir çok parçalı sistemi (mikroservisler, dağıtık sistemler, CI/CD pipeline'ları) öğrenirken de aynı sırayla ilerleyebilirsin. Takıldığın somut bir nokta olduğunda (bir hata mesajı, bir kavram, bir tasarım kararı), o noktayı anlatarak devam edebilirsin.