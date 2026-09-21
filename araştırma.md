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