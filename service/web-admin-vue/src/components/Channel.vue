<template>
    <div class="container p-5">
      <div class="text-center">
        <h1> {{this.$route.params.id}}</h1>
      </div>
    <div class="row row-cols-1 row-cols-md-4">
        <div class="col mb-4" v-for="v in msg" :key="v.agency_id">
            <div v-if="v.type == true" class="card border-success mb-3">
          <b-card border-variant="success" :header="v.agency_name" align="center">         
            <b-card-body>
              <b-card-title>Update setting</b-card-title>
                <div class="row row-cols-1">
                    <b-button @click="$bvModal.show('bv-modal-'+v.agency_name)" >Update setting</b-button> 
                    <b-modal :id="'bv-modal-'+v.agency_name" hide-footer>
                      <template #modal-title>
                        Update {{v.agency_name}}
                      </template>
                      <div class="d-block">
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
                            <input class="form-check-input" type="checkbox" name="lite_mode" v-model="v.lite_mode">
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
                        <div class="container p-1 pb-2">
                          <label class="form-label">Select channel type</label>
                          <div class="form-check">
                            <b-form-select v-model=v.type :options="v.channel_type"></b-form-select>
                          </div>
                        </div>
                      </div>        
                      <div class="container p-1">
                            <b-button>Update</b-button> 
                            <b-button>Disable agency</b-button>                         
                      </div>                             
                      </div>
                    </b-modal>                            
                </div>
            </b-card-body>
            </b-card>
            </div>

            <div v-if="v.type == false" class="card border-success mb-3">
            <h5 class="card-header">{{v.agency_name}}</h5>
            <div class="card-body">
                <div class="row row-cols-1">
                    <b-button @click="$bvModal.show('bv-modal-'+v.agency_name)" >Enable this agency</b-button> 
                    
                    <b-modal :id="'bv-modal-'+v.agency_name" hide-footer>
                      <template #modal-title>
                        Using <code>$bvModal</code> Methods
                      </template>
                      <div class="d-block text-center">
                        <h3>Hello From This Modal! {{v.agency_name}}</h3>
                      </div>
                    </b-modal>                            
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
      msg: [],
      form: {
        channel_type : [
          {value: 1,text: "Fanart"},
          {value: 2,text: "Livestream"},
          {value: 69,text: "Lewd"},
          {value: 3,text: "Fanart & Livestream"},
          {value: 70,text: "Fanart & Lewd"},
        ]
      }
    }
  },
  mounted(){
      axios.get("http://localhost:5000/channel/"+this.$route.params.id+"/agency",{
          withCredentials: true,
      })
      .then(response => {
        response.data.forEach((element) => {
          element.channel_type = [
          {value: 1,text: "Fanart"},
          {value: 2,text: "Livestream"},
          {value: 69,text: "Lewd"},
          {value: 3,text: "Fanart & Livestream"},
          {value: 70,text: "Fanart & Lewd"},
        ]
          this.msg.push(element)
        });
      })
  }
}
</script>