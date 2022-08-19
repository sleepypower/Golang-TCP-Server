<template>
  <div class="layoutContainer">
    <form v-on:submit.prevent="makeWebsiteThumbnail">
      <div class="form-group">
        <input
          v-model="websiteUrl"
          type="text"
          id="website-input"
          placeholder="Enter a website"
          class="form-control"
        />
      </div>
      <div class="form-group">
        <button class="btn btn-primary">Generate!</button>
      </div>
    </form>
    <div class="leftContainer">
      <p>GO TCP SERVER ANALYTICS</p>
      <h1>{{ users_connected }}</h1>
      <p>USERS CONNECTED</p>
      <hr class="rounded" />
    </div>
    <div class="rightContainer">
      <div class="topContainer">
        <h1>Users</h1>
        <!-- <vueper-slides>
          <vueper-slide v-for="i in 5" :key="i" :title="i.toString()" />
        </vueper-slides> -->
      </div>
      <div class="bottomContainer">
        <h1>Files</h1>
      </div>
    </div>
  </div>
</template>

<script>
// import { VueperSlides, VueperSlide } from "vueperslides";
// import "vueperslides/dist/vueperslides.css";
import axios from "axios";

export default {
  name: "HelloWorld",

  data() {
    return {
      users_connected: 0,
      files_sent: 0,
      bytes_sent: 0,
      channels: 0,
    };
  },

  methods: {
    makeWebsiteThumbnail() {
      // Call the Go API, in this case we only need the URL parameter.
      axios
        .post("http://localhost:3000/api/thumbnail", {})
        .then((response) => {
          console.log(response);
          console.log("hey!");
          this.users_connected = response.data.users_connected;
          this.files_sent = response.data.files_sent;
          this.bytes_sent = response.data.bytes_sent;
          this.channels = response.data.channels;
        })
        .catch((error) => {
          window.alert(`The API returned an error: ${error}`);
        });
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.layoutContainer {
  /* background-color: orange; */
  display: flex;
  min-height: 700px;
  padding: 20px 40px;
}

hr.rounded {
  border-top: 8px solid #ffbf0d;

  border-radius: 5px;
  width: 80%;
}

.leftContainer {
  background-color: #5100ff;
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  flex-grow: 0;
  margin: 20px 60px;
  border-radius: 8px;
  flex-shrink: 4;
  padding: 60px 20px;
}

.leftContainer p {
  font-size: 30px;
  text-align: center;
  font-weight: bold;
  color: white;
}

.leftContainer h1 {
  color: #ffbf0d;
  font-size: 150px;
  text-align: center;
  font-weight: bold;
}

.rightContainer {
  flex-shrink: 1;
  color: white;
  flex-grow: 7;
  margin: 20px 10px;
  display: flex;
  flex-direction: column;
}

.topContainer {
  flex-grow: 1;
  border-radius: 8px;
  background-color: red;
  padding: 0 20px;
  margin-bottom: 60px;
}

.topContainer h1 {
  text-align: left;
}

.bottomContainer {
  background-color: #5100ff;
  flex-grow: 1;
  background-color: 5100ff;
  border-radius: 8px;
}
.bottomContainer h1 {
  text-align: left;
  padding: 0 20px;
}
</style>
