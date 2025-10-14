// Глобальные переменные
let currentMod = null;

// Инициализация приложения
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
    setupEventListeners();
    loadAppInfo();
});

// Инициализация приложения
function initializeApp() {
    updateStatus('Готов к работе');
    console.log('HOI4 Visual Modder инициализирован');
}

// Настройка обработчиков событий
function setupEventListeners() {
    const selectModBtn = document.getElementById('selectModBtn');
    const modPathInput = document.getElementById('modPath');
    const editFocusBtn = document.getElementById('editFocusBtn');
    const editTechBtn = document.getElementById('editTechBtn');

    selectModBtn.addEventListener('click', handleSelectMod);
    modPathInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            handleSelectMod();
        }
    });

    editFocusBtn.addEventListener('click', () => {
        alert('Редактор фокусов будет реализован в следующей версии');
    });

    editTechBtn.addEventListener('click', () => {
        alert('Редактор технологий будет реализован в следующей версии');
    });
}

// Обработка выбора мода
async function handleSelectMod() {
    const modPath = document.getElementById('modPath').value.trim();
    
    if (!modPath) {
        showError('Пожалуйста, введите путь к каталогу мода');
        return;
    }

    updateStatus('Проверка каталога мода...');
    
    try {
        // Вызываем метод Go через Wails
        const result = await window.go.app.App.SelectModDirectory(modPath);
        
        if (result.isValid) {
            currentMod = result.modInfo;
            showModInfo(result);
            updateStatus(`Мод загружен: ${currentMod.name}`);
        } else {
            showValidationErrors(result);
            updateStatus('Ошибки валидации мода');
        }
    } catch (error) {
        console.error('Ошибка при выборе мода:', error);
        showError('Ошибка при проверке каталога мода: ' + error.message);
        updateStatus('Ошибка');
    }
}

// Отображение информации о моде
function showModInfo(result) {
    const modInfo = result.modInfo;
    
    // Показываем секцию информации о моде
    const modInfoSection = document.getElementById('modInfo');
    const modDetails = document.getElementById('modDetails');
    
    modDetails.innerHTML = `
        <div class="mod-summary">
            <h3>📋 ${modInfo.name}</h3>
            <p><strong>Путь:</strong> ${modInfo.basePath}</p>
            <p><strong>Статус:</strong> <span class="status-valid">✅ Валидный</span></p>
        </div>
        <div class="mod-stats">
            <div class="stat-item">
                <span class="stat-number">${modInfo.nationalFocusFiles.length}</span>
                <span class="stat-label">Файлов фокусов</span>
            </div>
            <div class="stat-item">
                <span class="stat-number">${modInfo.technologyFiles.length}</span>
                <span class="stat-label">Файлов технологий</span>
            </div>
        </div>
    `;
    
    modInfoSection.style.display = 'block';
    modInfoSection.classList.add('fade-in');
    
    // Показываем результаты валидации
    showValidationResults(result);
    
    // Показываем файлы
    showFiles(modInfo);
}

// Отображение результатов валидации
function showValidationResults(result) {
    const validationSection = document.getElementById('validationResults');
    const validationTitle = document.getElementById('validationTitle');
    const validationContent = document.getElementById('validationContent');
    
    let content = '';
    let titleClass = 'validation-success';
    let titleText = '✅ Проверка пройдена успешно';
    
    if (result.errors && result.errors.length > 0) {
        titleClass = 'validation-error';
        titleText = '❌ Обнаружены ошибки';
        
        content += '<div class="validation-group"><h4>Ошибки:</h4>';
        result.errors.forEach(error => {
            content += `<div class="validation-item validation-error">${error}</div>`;
        });
        content += '</div>';
    }
    
    if (result.warnings && result.warnings.length > 0) {
        if (titleClass === 'validation-success') {
            titleClass = 'validation-warning';
            titleText = '⚠️ Проверка пройдена с предупреждениями';
        }
        
        content += '<div class="validation-group"><h4>Предупреждения:</h4>';
        result.warnings.forEach(warning => {
            content += `<div class="validation-item validation-warning">${warning}</div>`;
        });
        content += '</div>';
    }
    
    if (!content) {
        content = '<div class="validation-item validation-success">Все проверки пройдены успешно!</div>';
    }
    
    validationTitle.textContent = titleText;
    validationTitle.className = titleClass;
    validationContent.innerHTML = content;
    
    validationSection.style.display = 'block';
    validationSection.classList.add('fade-in');
}

// Отображение ошибок валидации
function showValidationErrors(result) {
    showValidationResults(result);
    
    // Скрываем секции, которые не должны показываться при ошибках
    document.getElementById('modInfo').style.display = 'none';
    document.getElementById('filesSection').style.display = 'none';
}

// Отображение файлов
function showFiles(modInfo) {
    const filesSection = document.getElementById('filesSection');
    const focusFiles = document.getElementById('focusFiles');
    const techFiles = document.getElementById('techFiles');
    const editFocusBtn = document.getElementById('editFocusBtn');
    const editTechBtn = document.getElementById('editTechBtn');
    
    // Файлы фокусов
    if (modInfo.nationalFocusFiles.length > 0) {
        focusFiles.innerHTML = modInfo.nationalFocusFiles
            .map(file => `<div class="file-item">${file}</div>`)
            .join('');
        editFocusBtn.disabled = false;
    } else {
        focusFiles.innerHTML = '<div class="no-files">Файлы национальных фокусов не найдены</div>';
        editFocusBtn.disabled = true;
    }
    
    // Файлы технологий
    if (modInfo.technologyFiles.length > 0) {
        techFiles.innerHTML = modInfo.technologyFiles
            .map(file => `<div class="file-item">${file}</div>`)
            .join('');
        editTechBtn.disabled = false;
    } else {
        techFiles.innerHTML = '<div class="no-files">Файлы технологий не найдены</div>';
        editTechBtn.disabled = true;
    }
    
    filesSection.style.display = 'block';
    filesSection.classList.add('fade-in');
}

// Загрузка информации о приложении
async function loadAppInfo() {
    try {
        const appInfo = await window.go.app.App.GetAppInfo();
        document.getElementById('appInfo').textContent = 
            `${appInfo.name} v${appInfo.version} by ${appInfo.author}`;
    } catch (error) {
        console.error('Ошибка загрузки информации о приложении:', error);
    }
}

// Обновление статуса
function updateStatus(message) {
    document.getElementById('statusText').textContent = message;
}

// Показ ошибки
function showError(message) {
    alert('Ошибка: ' + message);
    updateStatus('Ошибка: ' + message);
}

// Утилиты для работы с DOM
function addClass(element, className) {
    if (element && !element.classList.contains(className)) {
        element.classList.add(className);
    }
}

function removeClass(element, className) {
    if (element && element.classList.contains(className)) {
        element.classList.remove(className);
    }
}
