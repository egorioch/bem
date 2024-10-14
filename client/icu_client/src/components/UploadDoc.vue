<template>
  <div>
    <h2>Upload Document</h2>
    <form @submit.prevent="submitForm">
      <div>
        <label for="name">Document Name:</label>
        <input type="text" v-model="meta.name" id="name" required />
      </div>

      <div>
        <label for="file">Select File:</label>
        <input type="file" @change="handleFileUpload" id="file" required />
      </div>

      <div>
        <label for="public">Public:</label>
        <input type="checkbox" v-model="meta.public" id="public" />
      </div>

      <div>
        <label for="grant">Grant Access to (comma separated logins):</label>
        <input type="text" v-model="grantList" id="grant" placeholder="login1, login2" />
      </div>

      <button type="submit">Upload</button>

      <!-- Показ результата загрузки -->
      <div v-if="statusMessage">
        <p>{{ statusMessage }}</p>
      </div>
    </form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      meta: {
        name: '',
        public: false,
        mime: '', // Автоматически определим MIME после выбора файла
        grant: [],
        file: false,
      },
      file: null,
      grantList: '',
      statusMessage: ''
    };
  },
  methods: {
    // Обработка загрузки файла
    handleFileUpload(event) {
      const file = event.target.files[0];
      this.meta.name = this.file.name
      this.file = file;
      this.meta.mime = file.type; // Устанавливаем MIME тип
      this.meta.file = true
    },

    // Отправка формы
    async submitForm() {
      // Преобразуем список пользователей с правами в массив
      this.meta.grant = this.grantList.split(',').map(login => login.trim());

      const formData = new FormData();

      // Добавляем метаданные в форму
      formData.append('meta', JSON.stringify(this.meta));

      // Добавляем файл
      if (this.file) {
        formData.append('file', this.file);
      } else {
        this.statusMessage = "Please select a file.";
        return;
      }

      try {
        // Отправляем POST-запрос на сервер
        const response = await axios.post('http://localhost:8080/api/docs/save', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
        this.statusMessage = response.data.status;
      } catch (error) {
        this.statusMessage = `Error: ${error.response?.data?.error || error.message}`;
      }
    }
  }
};
</script>

<style scoped>
/* Стили для формы */
form {
  display: flex;
  flex-direction: column;
}

label {
  margin-bottom: 0.5rem;
}

input {
  margin-bottom: 1rem;
}

button {
  margin-top: 1rem;
  padding: 0.5rem 1rem;
  cursor: pointer;
}

p {
  margin-top: 1rem;
  color: green;
}
</style>