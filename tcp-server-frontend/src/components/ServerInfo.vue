<template>
  <div class="serverInfoContainer">
    <!-- <h1>Server Info</h1> -->

    <div class="containerLeft">
      <h1>{{ bytes_sent }}</h1>
      <p>Bytes Sent</p>
      <hr class="rounded" />
    </div>
    <div class="containerRight">
      <div class="containerSmall">
        <p>
          Files Sent <span>{{ files_sent }}</span>
        </p>
      </div>
      <div class="containerSmall">
        <p>
          Channels <span>{{ channels }}</span>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
// import { VueperSlides, VueperSlide } from "vueperslides";
// import "vueperslides/dist/vueperslides.css";

export default {
  name: "HelloWorld",

  data() {
    return {};
  },
  props: {
    users_connected: {
      type: Number,
      default: 0,
    },
    files_sent: {
      type: Number,
      default: 0,
    },
    bytes_sent: {
      type: Number,
      default: 0,
    },
    channels: {
      type: Number,
      default: 0,
    },
  },

  methods: {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.serverInfoContainer {
  flex-grow: 1;
  display: flex;
  font-weight: 700;
  font-size: 25px;
}
.serverInfoContainer h1 {
  text-align: left;
  padding: 0 20px;
}

.containerLeft {
  background-color: #5100ff;
  color: white;
  display: flex;
  flex-direction: column;
  flex-grow: 0;
  justify-content: space-around;
  border-radius: 8px;

  padding: 30px 60px;
  max-width: 300px;
  margin-top: 20px;
}

.containerLeft h1 {
  color: #ffbf0d;
  font-size: 60px;
  text-align: center;
  font-weight: bold;
  padding: 0;
  margin: 0;
}

.containerLeft p {
  font-size: 30px;
  text-align: center;
  font-weight: bold;
  color: white;
  margin: 0;
}

hr.rounded {
  border-top: 8px solid #ffbf0d;

  border-radius: 5px;
  width: 40%;
}
/* ////////// */
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
  font-size: 200px;
  text-align: center;
  font-weight: bold;
}

/* ///////// */

.containerRight {
  flex-grow: 1;
  background-color: 5100ff;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  text-align: center;
  text-align: center;
}

.containerRight span {
  color: #ffbf0d;
  font-size: 40px;
  font-weight: 600;
}

.containerSmall {
  background-color: #5100ff;
  border-radius: 8px;
  height: 100%;
  margin-top: 20px;
  margin-left: 40px;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  align-items: center;
  padding: 0;
}
</style>
