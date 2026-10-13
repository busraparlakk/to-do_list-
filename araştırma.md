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