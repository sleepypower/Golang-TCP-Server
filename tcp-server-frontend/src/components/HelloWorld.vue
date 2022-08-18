<template>
  <div id="" class="container">
    <div class="row">
      <div class="col-md-6 offset-md-3 py-5">
        <h1>Generate a thumbnail of a website</h1>

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
        <h1>users_connected {{ users_connected }}</h1>
        <h1>files_sent {{ files_sent }}</h1>
        <h1>bytes_sent {{ bytes_sent }}</h1>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "HelloWorld",

  data() {
    return {
      users_connected: 0,
      files_sent: 0,
      bytes_sent: 0,
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
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
