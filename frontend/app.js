async function loadItems() {
    const ul = document.getElementById('items');
    ul.innerHTML = '<li>Загрузка...</li>';
    try {
        const res = await fetch('/api/items');
        if (!res.ok) throw new Error('HTTP ' + res.status);
        const items = await res.json();
        if (items.length === 0) {
            ul.innerHTML = '<li>Пусто</li>';
            return;
        }
        ul.innerHTML = '';
        items.forEach(it => {
            const li = document.createElement('li');
            li.textContent = `#${it.id} — ${it.name}`;
            ul.appendChild(li);
        });
    } catch (e) {
        ul.innerHTML = '<li style="color:red">Ошибка: ' + e.message + '</li>';
    }
}

document.getElementById('refresh').addEventListener('click', loadItems);
loadItems();