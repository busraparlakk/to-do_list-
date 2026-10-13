import express from "express";

const app = express();

app.listen(3000, () => {  // prtu dinle ve ekrana yaz 
    console.log("Express server 3000 portunda çalışıyor.");
});