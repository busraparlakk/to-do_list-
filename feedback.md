# 1: 

todo_handler.go: 
createTodoHandler fonksiyonu: 
```go
if err != nil {
		http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
		return
	}

	request.Title = strings.TrimSpace(request.Title)

	if request.Title == "" {
		http.Error(w, "title boş olamaz", http.StatusBadRequest)
		return
	}

	if request.DueDate != "" {
		_, err := time.Parse(time.RFC3339, request.DueDate)
		if err != nil {
			http.Error(w, "due_date geçerli bir RFC3339 tarihi olmalı", http.StatusBadRequest)
			return
		}
	}
```
Bu kod solid e uygun değil. Bu şekilde yeni bir obje kontrolü eklemek istediğinde aşağı bir if else daha yazarsın kod uzar büyür karmaşıklaşır. 
Bunun yerine todo.go dosyasında: 

```go 

type CheckResult struct{
    result boolean 
    desc string 
}

type Todo function(){

    check(){
        if request.Title == ""{
            return {false,"Title Empty"}
        }

        if request.DueDate......
    }
}
```
gibi bir kod yazarsan kontrollerini toplu yapmış olursun. ve createTodoHandler da şu satırı yazman yeterli olur: 

todo := Todo{
		ID:          hex.EncodeToString(idBytes),
		Title:       request.Title,
		Description: request.Description,
		Done:        false,
		DueDate:     request.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

if (!todo){
    return http.error(todo.result.desc)......
}

Aynı şekilde updateTodoHandler da aynı if bloklarını tekrar yazmışsın. Kod tekrarına düşmüş bu kontrolü eklersen kod tekrarı olmaz bir yerde yazarsın her yerde kullanırsın. kod bağımlılığını düşürerek kod karmaşıklığını ve spagetti kodu düşürebilirsin.

Bir yerde yaz güncellemen gerekirse tek bir yeri güncelle. SOLID in birinci kavramı. single responsibility. Herkes kendi görevini yapmalı updatetodohandler ın görevleri arasında obje kontrolü olmamalı mı gerçekten? 

go da bunu yapmanın çok daha basit bir yolu var hazır kütüphaneler ile yapabilirsin. Bununla beraber DTO mantığını da araştır. DTO nun ne olduğunu öğren. 

Ben kodları sözde yazdım syntax a bakmadım sadece anla diye. 

Bu



# 2: 

Gerçek dünyada veriler asla silinmez. Status alanı deleted olur: deleteTodoHandler fonksiyonu array siliyor. Burada veri silmeyi araştırabilirsin. İki en büyük yol var: 

Ya Deleted tablosu olur silinenler oraya taşınır ya da isDelete gibi alanlar olur. Veriler uygulamaya bağlı olarak 6 ay / 1 yıl sonra direkt silinir ya da silinmez. KVKK gibi durumlar burada devreye giriyor. Silinen veriler daha ne kadar ve nasıl saklanmalı vs. gibi. Bunu da araştırabilirsin: 

DELETE Best Practise? 