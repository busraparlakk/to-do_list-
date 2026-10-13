import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const DATA_FILE = path.join(__dirname, 'data', 'todos.json'); // dosya yolunu sakyalayak değişken oluşturur
function readTodos() {
    if (!fs.existsSync(DATA_FILE)) {  // dosyanın varlığını kontroleder 
     return [];
    }
    const metin = fs.readFileSync(DATA_FILE, 'utf8'); // okunan içeriği metin değişkeninde saklar
    if (metin.trim() === '') {
    return [];
    }
    return JSON.parse(metin);
}
function writeTodos(todos) {
    const metin = JSON.stringify(todos, null, 2); // javascript verisini jsona dönüştürür
    fs.writeFileSync(DATA_FILE, metin, 'utf8'); // dosyaya veri yazar 
}

export { readTodos, writeTodos };