<template>
<div class="container p-5">
    <div class="row row-cols-1 row-cols-md-4">
        <div class="col mb-4" v-for="v in msg.guild_channel" :key="v.id">
            <div class="card border-success mb-3">
            <h5 class="card-header">#{{v.name}}</h5>
            <div class="card-body">
                <div class="row row-cols-1">
                    <router-link :to="'/channel/'+v.id">
                        <b-button>Update channel</b-button>
                    </router-link>                  
                </div>
            </div>
            </div>   
        </div>       
    </div>
</div>    
</template>

<script>
import axios from 'axios';

export default {
  name: 'Channel',
  data () {
    return {
      msg: '',
    }
  },
  mounted(){
      axios.get("/v1/guilds/"+this.$route.params.id+"/channels",{
          withCredentials: true,
      })
      .then(response => {
          this.msg = response.data
      }).catch((error) => {
            this.$toasted.show("Oops somethings error "+error)
            console.log(error)
        }).finally(() => {
            //Perform action in always
        });
  },
}
</script>