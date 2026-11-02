# mikroservis mimarisi 
    bilinen monolit ve modüler mimarideki birleşik yapıdan farklı olarak sistemim her işlevi için bağımsız bir servis oluşturmayı amaçlayan mmimaridir. herhangi bir servisin çalışmaması ,bakım ya da güncellenmesi sırasında sistemin geri kalanı etkilenmez. 
# Reverse Proxy
    istemci sunucu arasına girer ve istekleri birden fazla sunucuya dağıtır. güenlik ve yönlendirmeyi amaçlar
# Load Balancer
     gelen istekleri birden fazla sunucuya dağıtır(balance denge ordan)
# API Gateway
    API isteklerinin geçiş kapısıdır. kimlik doğrulama,rate limit,routing (gelen isteğin ereye gideceğine karar verme işlemidir path'e bakarak) gibi işlemleri yönetir. amaç apileri merkezi yönetmek
Reverse Proxy isteği aracılara yönlendirir , Load Balancer yükü dağıtır, API Gateway ise API trafiğini yönetir.
# proxy etmek nedir
    bir isteğin iki taraf arasında doğrudan gitmesi yerine aradaki bir sunucu üzerinden iletilmesi demek
# endpoint nedir 
    normal testini yaptığımız endpoint . apinin belirli bir işi yaptığı adres örneğin GET / todos : bütün todoları ver 
# HTTP metotları ne yapıyor 
    bu uygulama için 4 temel metot kullanılacak normaş web dersindeki metotlar : GET (veri alma) , POST (yeni veri oluşturmak için) , PUT (mevcut veriyi güncelleme) , DELETE
# path nedir 
     url'nin hangi kaynağa eriştiğini gösterir
# request body
    clientin apiye gönderdiği veridir. nomal requestten farkı şuymuş : normal request paketin tamamını ifade eder bu ise paketin içindeki taşınan veri yükünü (payload) temsil eder.
# response body 
    apinin cliente geri gönderdiği veridir. normal responseden farkı paketin içindeki asıl veriyi taşımasıdır .

    server her zaman response gönderir
    client request    request:istemek st
# bazı önemli http status kodları
    200 ok :işşlwm başarılı
    201 created : yeni bir şey başarıyla oluşturuldu 
    204 no content : işlem başarılı ama response bdy yok
    400 bad request : cliet yanlış-eksik veri göndermiş
# önce en az bağımlılığı olan katmanı sağlamlaştır
    web uı apı olmadan test edilemez . ama apiyi çeşitli yçntemlerle test edebiliriz yani diğer kısımlar apiye bağlı ama api kimseye bağlı değil oyüzden onu sağlamlaştırarak başlamak mantıklı olur
# DTO (Data Transfer Object)
    istemcinin API’ye gönderebileceği verileri temsil eden yapıdır. Bizim projede Todo kaydın tamamını tutarken, CreateTodoRequest yalnızca oluşturma sırasında alınan title, description ve due_date alanlarını içerir. Böylece kullanıcıdan gelen bilgiler ile sunucunun yönettiği id, oluşturulma zamanı ve silinme zamanı gibi alanları ayırırız.

UpdateTodoRequest içinde pointer (*string, *bool) kullanmamızın nedeni, alanın gönderilmemesiyle sıfır değerinin gönderilmesini ayırt etmektir. Örneğin done gönderilmezse mevcut durumu koruruz; done: false gönderilirse görevi tamamlanmamış olarak güncelleriz. DTO veriyi taşır; başlığın boş olmaması ve tarihin geçerliliği gibi ortak kuralları ise Todo üzerindeki doğrulama metodu kontrol eder.

Veri silmede iki temel yaklaşım vardır: hard delete, kaydı kalıcı olarak kaldırır; soft delete, kaydı saklayıp silinmiş olarak işaretler. Biz önceden kaydı diziden çıkarıp dosyaya yazıyorduk. Şimdi deleted_at alanına silinme zamanını yazıyoruz. Alan null ise kayıt aktif, doluysa silinmiş kabul ediliyor.

Bu nedenle silinen todo dosyada kalıyor fakat listede görünmüyor, tekil olarak getirilemiyor ve güncellenemiyor. Bu yaklaşım ileride geri yükleme geliştirmemizi kolaylaştırır; şu anda geri yükleme özelliğimiz yok. Ayrıca soft delete, veriyi sonsuza kadar saklamak gerektiği anlamına gelmez. Kalıcı silmenin ne zaman yapılacağı, uygulamanın ihtiyacına göre ayrıca belirlenmelidir.
# DELETE Best Practices
“silme işlemi için iyi uygulamalar” anlamına gelir. Tek bir yöntem veya kütüphane adı değildir; silmenin nasıl güvenilir ve tutarlı tasarlanacağını anlatır.

Temel düşünce, verinin gerçekten kaldırılması mı yoksa silinmiş olarak işaretlenmesi mi gerektiğine karar vermektir. Geri yükleme ihtiyacı varsa soft delete yararlı olabilir. Kalıcı kaldırma gerekiyorsa hard delete kullanılır. “Gerçek uygulamalarda veri asla silinmez” diye genel bir kural yoktur.

Silmenin API davranışı da net olmalıdır. Bizim projede başarılı DELETE, gövdesiz 204 No Content döndürüyor. Silinmiş veya olmayan kayıtta 404 dönüyor. Silinen kayıt listelenmiyor ve güncellenemiyor.

Bir diğer kavram idempotency: aynı DELETE isteğini tekrar göndermek, sistemde yeni bir değişiklik oluşturmamalıdır. İlk isteğin 204, ikincinin 404 dönmesi buna aykırı değildir; önemli olan yanıtların aynı olması değil, kaydın silinmiş durumda kalmasıdır. Biz de tekrar DELETE gönderildiğinde silinme zamanının değişmediğini test ettik.

Uygulama büyüdüğünde, kullanıcının o kaydı silmeye yetkili olması, ilişkili kayıtların nasıl etkileneceği ve saklanan verilerin ne zaman kalıcı silineceği de tasarlanmalıdır. Biz şu anda soft delete davranışını ve tekrar edilen isteklerin tutarlılığını uyguladık.
# Authentication — Kimlik doğrulama
Authentication, isteği yapan kişinin iddia ettiği kullanıcı olduğunu doğrulama işlemidir. Örneğin kullanıcı e-posta ve parolasını gönderdiğinde sunucu bu bilgilerin doğruluğunu kontrol eder. Kimlik doğrulama yalnızca parola ile yapılmaz; passkey, güvenlik anahtarı veya ek doğrulama faktörleri de kullanılabilir. Amaç, kullanıcıyı güvenilir biçimde tanımaktır. OWASP Cheat Sheet Series

Todo projemize kullanıcı sistemi eklediğimizi düşünelim. Kullanıcının hesabına giriş yaparak kimliğini kanıtlaması authentication olur. Ancak bu, bütün kullanıcıların görevlerine erişebileceği anlamına gelmez.
# Authorization — Yetkilendirme
Authorization, kimliği belirlenen kullanıcının hangi kaynaklara erişebileceğini ve hangi işlemleri yapabileceğini kontrol eder. Bu kontrol yalnızca “admin mi, normal kullanıcı mı?” sorusundan ibaret değildir; kaydın kime ait olduğu da önemlidir. Yetkiler, ilgili her istekte sunucu tarafından kontrol edilmelidir. Arayüzde bir butonu gizlemek tek başına yetkilendirme sağlamaz. OWASP Cheat Sheet Series

Örneğin kullanıcı kendi todo kaydını silebilirken başka bir kullanıcının kaydını silememelidir. Giriş yapmış olması authentication kontrolünü geçtiğini gösterir; silmek istediği todo’nun kendisine ait olduğunun kontrol edilmesi ise authorization olur.
# Login — Sisteme giriş
Login, kullanıcının oturum açmak için geçtiği akıştır. Bir giriş formunu doldurmak veya “Google ile giriş yap” düğmesine basmak bu akışı başlatabilir. Sistem kimliği doğrular; başarılı olduğunda sonraki isteklerde kullanılabilecek bir oturum veya token oluşturabilir. Dolayısıyla login kullanıcıya görünen giriş süreciyken authentication bu süreçte yapılan kimlik doğrulamasıdır. OWASP Cheat Sheet Series

Kullanıcının her todo listeleme işleminde yeniden parola yazmaması için girişten sonraki isteklerle kullanıcı arasında bir bağlantı kurulması gerekir. Session ve token mekanizmaları burada devreye girer.
# Session — Oturum temelli yaklaşım
Yaygın sunucu taraflı session yaklaşımında, başarılı girişten sonra sunucu bir oturum kaydı oluşturur. Bu kaydı kullanıcıyla ilişkilendirir ve tarayıcıya tahmin edilmesi zor bir session ID verir. Bu kimlik genellikle cookie içinde tutulur. Tarayıcı uygun isteklere cookie’yi otomatik ekler; sunucu da session ID üzerinden kullanıcının oturumunu bulur. Cookie ile session aynı şey değildir: cookie kimliği taşır, session ise oturum bilgisini temsil eder. OWASP Cheat Sheet Series

Örneğin tarayıcı bir todo oluşturma isteği gönderdiğinde sunucu oturum kaydından kullanıcının kim olduğunu öğrenebilir. Çıkış yapıldığında sunucudaki oturum geçersiz kılınabilir. Birden fazla sunucu varsa oturum bilgilerinin bu sunucular arasında nasıl paylaşılacağı da tasarlanmalıdır.
# JWT — JSON Web Token
JWT, bilgileri belirli bir yapıda taşıyan bir token formatıdır. İçinde kullanıcı kimliği, token’ın kim tarafından üretildiği ve geçerlilik süresi gibi claim adı verilen alanlar bulunabilir. Yaygın kullanılan imzalı JWT üç parçadan oluşur: header, payload ve signature. Sunucu yalnızca token’ın içeriğini okumakla yetinmez; imzasını ve beklenen geçerlilik koşullarını da doğrulamalıdır. rfc-editor.org

Buradaki önemli ayrıntı, imzalı olmanın şifrelenmiş olmak anlamına gelmemesidir. Yaygın imzalı JWT’lerin payload bölümü okunabilir; imza, içeriğin değiştirilmesini tespit etmeye yarar. Bu nedenle token içine parola koyulmaz. JWT standardı şifrelenmiş yapıları da destekler, ancak her JWT şifreli değildir. rfc-editor.org

Todo projesinde JWT kullansaydık API, doğrulanmış token’daki kullanıcı kimliğini işlemi yapan kullanıcı olarak değerlendirebilirdi. Yine de ilgili todo’yu değiştirmeye yetkisi olup olmadığını ayrıca kontrol etmesi gerekirdi.
# Bearer — Token’ı sunma ve kullanma biçimi
Bearer, token’ı elinde bulunduran kişinin onun sağladığı erişimi kullanabilmesi anlamına gelir. Yaygın kullanımda istemci token’ı HTTP isteğinin şu başlığında gönderir: Authorization: Bearer <token>. Sunucu token’ı doğrular ve isteğin izin verilen kapsamda olup olmadığını kontrol eder. Token’ı ele geçiren biri de kullanabileceği için aktarımda HTTPS ile korunması gerekir. RFC Editor

Bearer ile JWT aynı şey değildir. JWT token’ın yapısını tanımlar; Bearer ise token’ın erişim için nasıl sunulduğunu ifade eder. Bearer olarak gönderilen token bir JWT olabilir veya sunucunun kayıtlarından kontrol ettiği, içeriği anlam taşımayan rastgele bir değer olabilir. Bu yüzden “JWT mi, Bearer mı?” şeklinde seçim yapmayız; bir JWT’yi Bearer olarak kullanabiliriz. rfc-editor.org

Bizim projeye uyarlarsak kullanıcı önce login sürecinden geçer, authentication ile kimliği doğrulanır ve sonraki isteklerde session veya token üzerinden tanınır. Ardından authorization kontrolü, örneğin silmek istediği todo’nun gerçekten kendisine ait olup olmadığını belirler. Gateway’in isteği doğru servise yönlendirmesi ise bu kontrolleri kendiliğinden sağlamaz; bunları ayrıca geliştirmemiz gerekir.