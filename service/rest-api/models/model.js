const { whereIn } = require("../config/db.js");
const knex = require("../config/db.js");


async function streamquery(Name,Table,Status,Limit){
  if (Name =="all"){
    return  await knex.select('VtuberMember.VtuberName','VtuberMember.VtuberName_EN',
    'VtuberMember.VtuberName_JP','VtuberGroup.VtuberGroupName',
    'VtuberMember.Youtube_ID','VtuberMember.Youtube_Avatar','VtuberMember.BiliBili_RoomID','VtuberMember.BiliBili_Avatar','VtuberMember.Region',
    ''+Table+'.*').from(Table)
  .innerJoin('VtuberMember',''+Table+'.VtuberMember_id','VtuberMember.id')
  .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
  .where(''+Table+'.Status',Status)
  .orderBy('ScheduledStart',"desc")
  .limit(Limit)

  } else {
    return await knex.select('VtuberMember.VtuberName','VtuberMember.VtuberName_EN',
    'VtuberMember.VtuberName_JP','VtuberGroup.VtuberGroupName',
    'VtuberMember.Youtube_ID','VtuberMember.BiliBili_RoomID','VtuberMember.BiliBili_Avatar','VtuberMember.Region',
    ''+Table+'.*').from(Table)
  .innerJoin('VtuberMember',''+Table+'.VtuberMember_id','VtuberMember.id')
  .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
  .where(function() {
    this.orWhereIn('VtuberMember.VtuberName', Name)
    .orWhereIn('VtuberMember.VtuberName_EN',Name)
    .orWhereIn('VtuberMember.VtuberName_JP',Name)
    .orWhereIn('VtuberGroup.VtuberGroupName',Name)
  })
  .andWhere(''+Table+'.Status',Status)
  .orderBy('ScheduledStart',"desc")
  .limit(Limit)
  }
}


const GetMemberAll = async (Reg,result) => {    
  try {
    let data
    if(Reg != null) {
      data = await knex.select('VtuberMember.*','VtuberGroup.VtuberGroupName').from('VtuberMember')
      .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
      .where('Region',Reg)
      .orderBy(['VtuberGroup_id','Region'])        
    } else {
      data = await knex.select('VtuberMember.*','VtuberGroup.VtuberGroupName').from('VtuberMember')
      .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
      .orderBy(['VtuberGroup_id','Region'])    
    }
    if (data.length){
      data.forEach(i => {
        delete i.id
        delete i.VtuberGroup_id
      });
      result(null,data)
    } else {
      result({ kind: "not_found" }, null);
    }
  } catch (error) {
    console.log(error)      
    result({kind:"Error kntl"},null)
  }
};


const GetMemberName = async (Name,Region, result) => {
  try {
    let data
    if (Region != null) {
      data = await knex('VtuberMember')
      .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
      .where(function() {
        this.orWhereIn('VtuberMember.VtuberName', Name)
        .orWhereIn('VtuberMember.VtuberName_EN',Name)
      }).orWhereIn('VtuberGroup.VtuberGroupName',Name)
      .andWhere(function(){
        this.whereIn("Region",Region)
      })

    } else{
      data = await knex('VtuberMember')
      .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
      .where(function() {
        this.orWhereIn('VtuberMember.VtuberName', Name)
        .orWhereIn('VtuberMember.VtuberName_EN',Name)
      }).orWhereIn('VtuberGroup.VtuberGroupName',Name)
      .orderBy("Region","desc")
    }

    if (data.length){
      data.forEach(i => {
        delete i.id
        delete i.VtuberGroup_id
      });
      result(null,data)
    } else {
      result({ kind: "not_found" }, null);
    }
  } catch (error) {
    result({kind:"Error cok"},null)
  }
};

const GetGroupAll = async result => {    
  try {
    let data = await knex('VtuberGroup')
    if (data.length){
      data.forEach(i => {
        delete i.id
      });
      result(null,data)
    } else {
      result({ kind: "not_found" }, null);
    }

  } catch (error) {
    console.log(error)      
    result({kind:"Error kntl"},null)
  }
};


const GetGroupName = async (Name, result) => {
  try {
    let data = await knex('VtuberGroup').whereIn('VtuberGroupName',Name)
    if (data.length){
      data.forEach(i => {
        delete i.id
      });
      result(null,data)
    } else {
      result({ kind: "not_found" }, null);
    }
  } catch (error) {
    console.log(error)
    result({kind:"Error kntl"},null)
  }
};

const GetTwitter = async (Name, Limit,result) => {
  try {
    let data = await knex.select('VtuberMember.VtuberName','VtuberMember.VtuberName_EN',
    'VtuberMember.VtuberName_JP','VtuberMember.Hashtag','VtuberGroup.VtuberGroupName',
    'Twitter.*').from('Twitter')
  .innerJoin('VtuberMember','Twitter.VtuberMember_id','VtuberMember.id')
  .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
  .where(function() {
    this.orWhereIn('VtuberMember.VtuberName', Name)
    .orWhereIn('VtuberMember.VtuberName_EN',Name)
    .orWhereIn('VtuberMember.VtuberName_JP',Name)
  }).orWhereIn('VtuberGroup.VtuberGroupName',Name)
  .orderBy("Twitter.id",'desc')
  .limit(Limit)
    if (data.length){
      data.forEach(i => {
        delete i.id
        delete i.VtuberMember_id
        i.Photos = i.Photos.split("\n");
      });
      result(null,data)
    } else {
      result({ kind: "not_found" }, null);
    }
  } catch (error) {
    console.log(error)
    result({kind:"Error kntl"},null)
  }
};


const GetTBilibili = async (Name, Limit,result) => {
  try {
    let data = await knex.select('VtuberMember.VtuberName','VtuberMember.VtuberName_EN',
    'VtuberMember.VtuberName_JP','VtuberMember.BiliBili_Hashtag','VtuberGroup.VtuberGroupName',
    'TBiliBili.*').from('TBiliBili')
  .innerJoin('VtuberMember','TBiliBili.VtuberMember_id','VtuberMember.id')
  .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
  .where(function() {
    this.orWhereIn('VtuberMember.VtuberName', Name)
    .orWhereIn('VtuberMember.VtuberName_EN',Name)
    .orWhereIn('VtuberMember.VtuberName_JP',Name)
  }).orWhereIn('VtuberGroup.VtuberGroupName',Name)
  .orderBy("TBiliBili.id",'desc')
  .limit(Limit)

    if (data.length){
      data.forEach(i => {
        delete i.id
        delete i.VtuberMember_id
        i.Photos = i.Photos.split("\n");
      });
      result(null,data)
    } else {
      result({ kind: "not_found" }, null);
    }
  } catch (error) {
    console.log(error)
    result({kind:"Error kntl"},null)
  }
};



const GetYtLivestream = async (Name, Status, Limit,result) => {
  console.log(Name, Status, Limit)
  try {
    data = await streamquery(Name,'Youtube',Status,Limit)
    if (data.length){
      data.forEach(i => {
        delete i.id
        delete i.VtuberGroup_id
        delete i.VtuberMember_id
        delete i.BiliBili_SpaceID
        delete i.BiliBili_Avatar
        i.Viewers = parseInt(i.Viewers,10)
      });
      result(null,data)
    } else {
      result({ kind: "not_found" }, null);
    }
  } catch (error) {
    console.log(error)
    result({kind:"Error kntl"},null)
  }
};



const GetLiveBilibili = async (Name, Status, Limit,result) => {
  try {
    data = await streamquery(Name,'LiveBiliBili',Status,Limit)
    if (data.length){
      data.forEach(i => {
        delete i.id
        delete i.VtuberGroup_id
        delete i.VtuberMember_id
        delete i.Youtube_ID
        delete i.Youtube_Avatar
      });
      result(null,data)
    } else {
      result({ kind: "not_found" }, null);
    }
  } catch (error) {
    console.log(error)
    result({kind:"Error kntl"},null)
  }
};


const GetSpaceBiliBili = async (Name, Limit,result) => {
  try {
    let data = await knex.select('VtuberMember.VtuberName','VtuberMember.VtuberName_EN',
    'VtuberMember.VtuberName_JP','VtuberGroup.VtuberGroupName',
    'VtuberMember.Youtube_ID','VtuberMember.BiliBili_SpaceID','VtuberMember.Region',
    'BiliBili.*').from('BiliBili')
  .innerJoin('VtuberMember','BiliBili.VtuberMember_id','VtuberMember.id')
  .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
  .where(function() {
    this.orWhereIn('VtuberMember.VtuberName', Name)
    .orWhereIn('VtuberMember.VtuberName_EN',Name)
    .orWhereIn('VtuberMember.VtuberName_JP',Name)
    .orWhereIn('VtuberGroup.VtuberGroupName',Name)
  })
  .orderBy('UploadDate',"desc")
  .limit(Limit)
  if(data.length){
    data.forEach(i => {
      delete i.id
      delete i.VtuberGroup_id
      delete i.VtuberMember_id
      i.Viewers = parseInt(i.Viewers,10)
    });
  
    result(null,data) 
  } else {
    result({ kind: "not_found" }, null); 
  }

  } catch (error) {
    console.log(error)
    result({kind:"Error kntl"},null)
  }
};


const Getsubscriber = async (Name, result) => {
  try {
    let data = await knex.select('VtuberMember.VtuberName','VtuberMember.VtuberName_EN',
    'VtuberMember.VtuberName_JP','VtuberGroup.VtuberGroupName',
    'VtuberMember.Youtube_ID','VtuberMember.BiliBili_SpaceID','VtuberMember.Region',
    'Subscriber.*').from('Subscriber')
  .innerJoin('VtuberMember','Subscriber.VtuberMember_id','VtuberMember.id')
  .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
  .where(function() {
    this.orWhereIn('VtuberMember.VtuberName', Name)
    .orWhereIn('VtuberMember.VtuberName_EN',Name)
    .orWhereIn('VtuberMember.VtuberName_JP',Name)
    .orWhereIn('VtuberGroup.VtuberGroupName',Name)
  })
  .orderBy('Region')
  if (data.length) {
    data.forEach(i => {
      delete i.id
      delete i.VtuberGroup_id
      delete i.VtuberMember_id
    });
    result(null,data)

  } else {
    result({ kind: "not_found" }, null); 
  }

  } catch (error) {
    console.log(error)
    result({kind:"Error kntl"},null)
  }
};


const GetDiscordChannel = async (ID,result) =>{
  try {
    let data = await knex.select('Channel.*','VtuberGroup.VtuberGroupName')
    .from("Channel").innerJoin('VtuberGroup','Channel.VtuberGroup_id','VtuberGroup.id')
    .where("Channel.DiscordChannelID",ID)
    Datafix = {
      "DiscordChannelID" : data[0].DiscordChannelID,
      "ChannelData" : []
    }
    if (data.length) {
      data.forEach(i => {
        if (i.Type == 3) {
          i.Type = "All"
        } else if (i.Type == 2) {
          i.Type = "Live"
        } else {
          i.Type = "Art"
        }
  
        Datafix.ChannelData.push({
          "id": i.id,
          "GroupName": i.VtuberGroupName,
          "Type": i.Type,
          "LiveOnly": Boolean(i.LiveOnly),
          "NewUpcoming": Boolean(i.NewUpcoming),
        })
      });
      result(null,Datafix)
    } else {
      result({ kind: "not_found" }, null); 
    }
  } catch (error) {
    console.log(error)
    result({kind:"Error kntl"},null)
  }
}

const GetRols = async (ChanndlID,result) =>{
  try {
    let data = await knex.select('Channel.*','VtuberGroup.VtuberGroupName')
    .from("Channel").innerJoin('VtuberGroup','Channel.VtuberGroup_id','VtuberGroup.id')
    .where("Channel.DiscordChannelID",ID)
    Datafix = {
      "DiscordChannelID" : data[0].DiscordChannelID,
      "ChannelData" : []
    }
    if (data.length) {
      data.forEach(i => {
        if (i.Type == 3) {
          i.Type = "All"
        } else if (i.Type == 2) {
          i.Type = "Live"
        } else {
          i.Type = "Art"
        }
  
        Datafix.ChannelData.push({
          "id": i.id,
          "GroupName": i.VtuberGroupName,
          "Type": i.Type,
          "LiveOnly": Boolean(i.LiveOnly),
          "NewUpcoming": Boolean(i.NewUpcoming),
        })
      });
      result(null,Datafix)
    } else {
      result({ kind: "not_found" }, null); 
    }
  } catch (error) {
    console.log(error)
    result({kind:"Error kntl"},null)
  }
}

module.exports = {
  GetMemberAll: GetMemberAll,
  GetMemberName: GetMemberName,
  GetGroupName : GetGroupName,
  GetGroupAll: GetGroupAll,
  GetYtLivestream : GetYtLivestream,
  GetTwitter : GetTwitter,
  GetTBilibili : GetTBilibili,
  GetLiveBilibili: GetLiveBilibili,
  GetSpaceBiliBili: GetSpaceBiliBili,
  Getsubscriber: Getsubscriber,
  GetDiscordChannel: GetDiscordChannel,
}