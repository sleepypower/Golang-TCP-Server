<template>
  <div class="layoutContainer">
    <!-- <form v-on:submit.prevent="makeWebsiteThumbnail">
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
    </form> -->
    <div class="leftContainer">
      <p class="yellowColor">GO TCP SERVER ANALYTICS</p>
      <h3>{{ users_connected }}</h3>
      <p>USERS CONNECTED</p>
      <hr class="rounded" />
      <button class="requestButton" @click="makeWebsiteThumbnail">
        Refresh
      </button>
    </div>
    <div class="rightContainer">
      <div class="topContainer">
        <h1>Users</h1>
        <vueper-slides
          class="no-shadow"
          :visible-slides="3"
          slide-multiple
          :gap="3"
          :slide-ratio="1 / 4"
          :dragging-distance="200"
          :breakpoints="{ 800: { visibleSlides: 2, slideMultiple: 2 } }"
        >
          <vueper-slide
            v-for="(slide, i) in clients"
            :key="i"
            :title="slide.username"
            :content="slide.channels"
            class="slide"
          />
        </vueper-slides>
      </div>
      <ServerInfo
        :bytes_sent="bytes_sent"
        :channels="channels"
        :files_sent="files_sent"
        :users_connected="users_connected"
      />
    </div>
  </div>
</template>

<script>
import { VueperSlides, VueperSlide } from "vueperslides";
import "vueperslides/dist/vueperslides.css";
import axios from "axios";
import ServerInfo from "./ServerInfo.vue";

export default {
  name: "HelloWorld",
  components: { ServerInfo, VueperSlides, VueperSlide },
  data() {
    return {
      users_connected: 0,
      files_sent: 0,
      bytes_sent: 0,
      channels: 0,
      clients: [],
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
          this.clients = response.data.clients;
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
.slideClient {
  background-color: wheat;
}

.layoutContainer {
  /* background-color: orange; */
  display: flex;
  height: 80vh;
  padding: 80px 120px;
}

hr.rounded {
  border-top: 8px solid #ffbf0d;
  border-radius: 5px;
  width: 80%;
}

.yellowColor {
  color: #ffbf0d !important;
}

.requestButton {
  background-color: #ffbf0d !important;
  border-top: 0px solid #ffbf0d;
  border-radius: 8px;
  margin: 0 auto;
  color: white;
  padding: 10px 20px;
  width: 60%;
  font-size: 1.2rem;
}

.slide {
  background-color: #5100ff;
  border-radius: 8px;
  border: 0px solid #5100ff;
  font-weight: 700;
}

.leftContainer {
  background-color: #5100ff;
  color: white;
  display: flex;
  flex-direction: column;
  flex-grow: 0;
  justify-content: space-around;
  margin-right: 80px;
  border-radius: 8px;
  flex-shrink: 4;
  padding: 60px;
  max-width: 300px;
}

.leftContainer p {
  font-size: 30px;
  text-align: center;
  font-weight: bold;
  color: white;
  margin: 0;
}

.leftContainer h3 {
  color: #ffbf0d;
  font-size: 150px;
  text-align: center;
  font-weight: bold;
  padding: 0;
  margin: 0;
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
  flex-grow: 2;
  border-radius: 8px;
  /* background-color: red; */
  padding: 0 20px;
  /* margin-bottom: 60px; */
}

.topContainer h1 {
  text-align: left;
}
</style>
