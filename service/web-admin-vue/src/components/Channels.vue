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

<!--
            <b-modal :id="'bv-modal-'+v.id" hide-footer>
                <template slot="modal-title">
                Update channel #{{v.name}}
                </template>
                <div class="d-block">
                  <select class="form-select select-agency" @change="select_agency($event)" :id="v.id" aria-label="Select Agency">
                      <option selected>Select Agency</option>
                      <option v-for="vv in msg.agency" :key=vv.id :value=vv.id > {{vv.VtuberGroupName}}</option>
                  </select>
                    <form action="#" :id="v.id+'-form'" method="post">
                      <div class="container p-1">
                        <div class="form-check form-switch">
                          <input class="form-check-input" type="checkbox" name="channel-status" v-model="form.channel_status" :id="'channel-status-'+v.id">
                          <label class="form-check-label" for="flexSwitchCheckDefault">Channel status </label>
                        </div>
                      </div>
                      <div class="container p-1">
                        <label class="form-label">Channel Option</label>                         
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="dynamic" v-model="form.dynamic" :id="'dynamic-'+v.id">
                            <label class="form-check-label" for="checkbox">
                                dynamic mode
                            </label>
                            <b-icon v-b-popover.hover.top="'Set channel to dynamic notification,Remove past livestream notification from discord channel'" title="Dynamic" icon="info-circle"></b-icon>
                        </div>
                        <div class="form-check">
                          <input class="form-check-input" type="checkbox" name="upcoming" v-model="form.upcoming" :id="'upcoming-'+v.id">
                          <label class="form-check-label" for="checkbox">
                              upcoming
                          </label>
                          <b-icon v-b-popover.hover.top="'Enable new livestreams notification'" title="Upcoming" icon="info-circle"></b-icon>
                        </div>                        
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="lite_mode" v-model="form.lite_mode" id="lite_mode">
                            <label class="form-check-label" for="checkbox">
                                lite mode 
                            </label>
                            <b-icon v-b-popover.hover.top="'Disabling ping user/role function'" title="Lite mode" icon="info-circle"></b-icon>
                        </div>                      
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="indie_notif" v-model="form.indie_notif" :id="'indie_notif-'+v.id" :disabled='form.disable_indie'>
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
                          v-model="form.channel_region"
                          :options="all_region"
                          class="mb-3"
                          value-field="item"
                          text-field="item"
                        ></b-form-checkbox-group>
                        </div>
                        <div class="container p-1 pb-2">
                          <select class="form-select channel-type" aria-label="select channel type" :id="'channel-type-'+v.id" name="channel">
                              <option value="-1" id="null">Select channel type</option>
                              <option value="1" :id="'fanart-'+v.id">Fanart</option>
                              <option value="2" :id="'livestream-'+v.id">Livestream</option>
                              <option value="69" :id="'lewd-'+v.id">Lewd</option>
                              <option value="3" :id="'fanart-livestream-'+v.id">Fanart & Livestream</option>
                              <option value="70" :id="'fanart-lewd-'+v.id">Fanart & Lewd</option>
                          </select>
                        </div>

                        <button type="submit" :id="'submit-'+v.id" value="Submit" class="btn btn-primary">Submit</button>
                      </div>
                    </form>                                  
                </div>
                <b-button class="mt-3" block @click="$bvModal.hide('bv-modal-example')">Close Me</b-button>
            </b-modal>    
-->         
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
      axios.get("http://localhost:5000/guilds/"+this.$route.params.id+"/channels",{
          withCredentials: true,
      })
      .then(response => {
          this.msg = response.data
      })
  },
}
</script>