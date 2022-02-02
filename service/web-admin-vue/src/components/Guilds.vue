<template>
    <div class="container p-5">
        <div class="row row-cols-1 row-cols-md-4">
            <div class="col mb-4" v-for="v in msg['guilds']" :key="v.id">
                <div class="card" style="width: 18rem;">
                    <img :src="'https://cdn.discordapp.com/icons/'+v.id+'/'+v.icon+'.png?size=4096'" class="card-img-top">
                    <div class="card-body">
                        <h5 class="card-title">{{v.name}}</h5>
                        <router-link :to="'/guilds/'+v.id+'/channels'">
                            <button class="btn btn-primary">Settings</button>
                        </router-link>                        
                    </div>
                </div>                     
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'Guilds',
  data () {
    return {
      msg: ''
    }
  },
  mounted(){
      axios.get("/v1/guilds",{
          withCredentials: true,
      })
      .then(response => {
          this.msg = response.data
      }).catch((error) => {
        this.$toasted.show("Oops somethings error "+error);
        console.log(error)
    }).finally(() => {
        //Perform action in always
    });
  }
}
</script>