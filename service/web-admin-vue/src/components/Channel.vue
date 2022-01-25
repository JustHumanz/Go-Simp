<template>
    <div class="container p-5">
      <div class="text-center">
        <h1> {{channel_name}}</h1>
      </div>
    <div class="row row-cols-1 row-cols-md-4">
        <div class="col mb-4" v-for="v in msg" :key="v.agency_id">
            <div v-if="v.type != null" class="card border-success mb-3">
          <b-card border-variant="success" :header="v.agency_name" align="center">         
            <b-card-body>
                <div class="row row-cols-1">
                    <b-button @click="v.modal = !v.modal" >Update setting</b-button> 
                    
                    <b-modal v-model="v.modal" hide-footer>
                      <template #modal-title>
                        Update {{v.agency_name}}
                      </template>
                      <div class="d-block">
                      <form v-on:submit.prevent="submitForm(v.agency_id)">
                      <div class="container p-1">
                        <label class="form-label">Channel Option</label>                         
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="dynamic" v-model="v.dynamic">
                            <label class="form-check-label" for="checkbox">
                                dynamic mode 
                            </label>
                            <b-icon v-b-popover.hover.top="'Set channel to dynamic notification,Remove past livestream notification from discord channel'" title="Dynamic" icon="info-circle"></b-icon>
                        </div>
                        <div class="form-check">
                          <input class="form-check-input" type="checkbox" name="upcoming" v-model="v.upcoming">
                          <label class="form-check-label" for="checkbox">
                              upcoming
                          </label>
                          <b-icon v-b-popover.hover.top="'Enable new livestreams notification'" title="Upcoming" icon="info-circle"></b-icon>
                        </div>                        
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="lite_mode" v-model="v.lite">
                            <label class="form-check-label" for="checkbox">
                                lite mode 
                            </label>
                            <b-icon v-b-popover.hover.top="'Disabling ping user/role function'" title="Lite mode" icon="info-circle"></b-icon>
                        </div>                      
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="indie_notif" v-model="v.indie_notif" :disabled='v.agency_id != "10"'>
                            <label class="form-check-label" for="checkbox">
                                independent notif
                            </label>
                            <b-icon v-b-popover.hover.top="'Send all independent vtubers notification'" title="independent notif" icon="info-circle"></b-icon>
                        </div>  
                      </div>        
                      <div class="container p-1">
                        <label class="form-label">Region</label>
                        <div class="form-check">
                        <b-form-checkbox-group
                          v-model="v.channel_region"
                          :options="v.agency_region"
                          class="mb-3"
                          value-field="item"
                          text-field="item"
                        ></b-form-checkbox-group>
                        </div>
                      </div>        
                      <div class="container p-1">
                        <label class="form-label">Select channel type</label>
                        <div class="form-check">
                          <select class="form-select"  v-model="v.type" aria-label="select channel type" >
                            <option v-for="option in v.channel_type" :value="option.text" :key="option.text">
                              {{ option.text }}
                            </option>
                          </select>
                        </div>
                      </div>                                
                      <div class="container p-1">
                        <div class="form-group">
                            <b-button type="submit">Update</b-button> 
                        </div>                                                 
                      </div>               
                      <div class="container p-1">
                        <b-button variant="danger">Disable agency</b-button>
                      </div>
                      </form>                          
                      </div>
                    </b-modal>                            
                </div>
            </b-card-body>
            </b-card>
            </div>

            <b-card v-if="v.type == null" border-variant="warning" :header="v.agency_name" align="center">  
            <div class="card-body">
                <div class="row row-cols-1">
                    <b-button @click="v.modal = !v.modal" >Enable this agency</b-button> 
                    
                    <b-modal v-model="v.modal" hide-footer>
                      <template #modal-title>
                        Enabeling {{v.agency_name}}
                      </template>
                      <div class="d-block">
                      <form v-on:submit.prevent="submitForm(v.agency_id)">
                      <div class="container p-1">
                        <label class="form-label">Channel Option</label>                         
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="dynamic" v-model="v.dynamic">
                            <label class="form-check-label" for="checkbox">
                                dynamic mode 
                            </label>
                            <b-icon v-b-popover.hover.top="'Set channel to dynamic notification,Remove past livestream notification from discord channel'" title="Dynamic" icon="info-circle"></b-icon>
                        </div>
                        <div class="form-check">
                          <input class="form-check-input" type="checkbox" name="upcoming" v-model="v.upcoming">
                          <label class="form-check-label" for="checkbox">
                              upcoming
                          </label>
                          <b-icon v-b-popover.hover.top="'Enable new livestreams notification'" title="Upcoming" icon="info-circle"></b-icon>
                        </div>                        
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="lite_mode" v-model="v.lite">
                            <label class="form-check-label" for="checkbox">
                                lite mode 
                            </label>
                            <b-icon v-b-popover.hover.top="'Disabling ping user/role function'" title="Lite mode" icon="info-circle"></b-icon>
                        </div>                      
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="indie_notif" v-model="v.indie_notif" :disabled='v.agency_id != "10"'>
                            <label class="form-check-label" for="checkbox">
                                independent notif
                            </label>
                            <b-icon v-b-popover.hover.top="'Send all independent vtubers notification'" title="independent notif" icon="info-circle"></b-icon>
                        </div>  
                      </div>        
                      <div class="container p-1">
                        <label class="form-label">Region</label>
                        <div class="form-check">
                        <b-form-checkbox-group
                          v-model="v.channel_region"
                          :options="v.agency_region"
                          class="mb-3"
                          value-field="item"
                          text-field="item"
                        ></b-form-checkbox-group>
                        </div>
                      </div>        
                      <div class="container p-1">
                        <label class="form-label">Select channel type</label>
                        <div class="form-check">
                          <select class="form-select"  v-model="v.type" aria-label="select channel type" >
                            <option v-for="option in v.channel_type" :value="option.text" :key="option.text">
                              {{ option.text }}
                            </option>
                          </select>
                        </div>
                      </div>                                
                      <div class="container p-1">
                        <div class="form-group">
                            <b-button type="submit">Update</b-button> 
                        </div>                                                 
                      </div>               
                      <div class="container p-1">
                        <b-button variant="danger">Disable agency</b-button>
                      </div>
                      </form>                          
                      </div>                      
                    </b-modal>                           
                </div>
            </div>
            </b-card>       
        </div>
    </div>
    </div>
</template>

<script>
import axios from 'axios';
import Config from "../config.json";

export default {
  name: 'Channel',
  data () {
    return {
      msg: [],
      channel_name: '',
      is_nsfw: '',
    }
  },
  mounted(){
      axios.get(Config.REST_API+"/channel/"+this.$route.params.id+"/agency",{
          withCredentials: true,
      })
      .then(response => {
        this.channel_name = response.data['channel_name']
        this.is_nsfw = response.data['is_nsfw']

        response.data['agency_list'].forEach((element) => {
          if (response.data['is_nsfw']) {
          element.channel_type = [
            {text: "Fanart"},
            {text: "Livestream"},
            {text: "Lewd"},
            {text: "Fanart & Livestream"},
            {text: "Fanart & Lewd"},
          ]
          } else{
          element.channel_type = [
            {text: "Fanart"},
            {text: "Livestream"},
            {text: "Fanart & Livestream"},
          ]            
          }
        element.modal = false
          this.msg.push(element)
        });
      })
  }, methods:{
      submitForm(agency_id){
        this.msg.forEach((element) => {
          if (element.agency_id == agency_id) {            
            axios.post(Config.REST_API+'/channel/'+this.$route.params.id+'/update', element,{withCredentials: true,})
                  .then((res) => {
                      console.log(res)
                  })
                  .catch((error) => {
                      // error.response.status Check status code
                      console.log(error)
                  }).finally(() => {
                      //Perform action in always
                  });                    
          }
        });
      }
  }
}
</script>