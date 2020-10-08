const knex = require("../config/db.js");


async function streamquery(Name,Table,Status,Limit){
  let data = await knex.select('VtuberMember.VtuberName','VtuberMember.VtuberName_EN',
  'VtuberMember.VtuberName_JP','VtuberGroup.VtuberGroupName',
  'VtuberMember.Youtube_ID','VtuberMember.BiliBili_SpaceID','VtuberMember.Region',
  ''+Table+'.*').from(Table)
.innerJoin('VtuberMember',''+Table+'.VtuberMember_id','VtuberMember.id')
.innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
.where(function() {
  this.where('VtuberMember.VtuberName', Name)
  .orWhere('VtuberMember.VtuberName_EN',Name)
  .orWhere('VtuberMember.VtuberName_JP',Name)
  .orWhere('VtuberGroup.VtuberGroupName',Name)
})
.andWhere(''+Table+'.Status',Status)
.orderBy('ScheduledStart',"desc")
.limit(Limit)
return data
}


const GetMemberAll = async result => {    
  try {
    let data = await knex('VtuberMember')
    console.log(data)
    data.forEach(i => {
      delete i.id
      delete i.VtuberGroup_id
    });
    result(null,data)
  } catch (error) {
    console.log(error)      
    result({kind:"Error kntl"},null)
  }
};


const GetMemberName = async (Name, result) => {
  try {
    let data = await knex('VtuberMember').whereIn('VtuberName_EN',Name)
      .orWhereIn('VtuberName',Name)
      .orWhereIn('VtuberName_JP',Name)
    if (data.length){
      console.log(data)
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
    result({kind:"Error cok"},null)
  }
};


const GetGroupAll = async result => {    
  try {
    let data = await knex('VtuberGroup')
    console.log(data)
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
      console.log(data)
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
    this.where('VtuberMember.VtuberName', Name)
    .orWhere('VtuberMember.VtuberName_EN',Name)
    .orWhere('VtuberMember.VtuberName_JP',Name)
  }).orWhere('VtuberGroup.VtuberGroupName',Name)
  .orderBy("Twitter.id",'desc')
  .limit(Limit)

    if (data.length){
      console.log(data)
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
    this.where('VtuberMember.VtuberName', Name)
    .orWhere('VtuberMember.VtuberName_EN',Name)
    .orWhere('VtuberMember.VtuberName_JP',Name)
  }).orWhere('VtuberGroup.VtuberGroupName',Name)
  .orderBy("TBiliBili.id",'desc')
  .limit(Limit)

    if (data.length){
      console.log(data)
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
    console.log(data)
    if (data.length){
      console.log(data)
      data.forEach(i => {
        i.Youtube_ID = i.Youtube_ID.split("\n")
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



const GetLiveBilibili = async (Name, Status, Limit,result) => {
  console.log(Name, Status, Limit)
  try {
    data = await streamquery(Name,'LiveBiliBili',Status,Limit)
    if (data.length){
      console.log(data)
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


const GetSpaceBiliBIli = async (Name, Limit,result) => {
  console.log(Name, Limit)
  try {
    let data = await knex.select('VtuberMember.VtuberName','VtuberMember.VtuberName_EN',
    'VtuberMember.VtuberName_JP','VtuberGroup.VtuberGroupName',
    'VtuberMember.Youtube_ID','VtuberMember.BiliBili_SpaceID','VtuberMember.Region',
    'BiliBili.*').from('BiliBili')
  .innerJoin('VtuberMember','BiliBili.VtuberMember_id','VtuberMember.id')
  .innerJoin('VtuberGroup','VtuberMember.VtuberGroup_id','VtuberGroup.id')
  .where(function() {
    this.where('VtuberMember.VtuberName', Name)
    .orWhere('VtuberMember.VtuberName_EN',Name)
    .orWhere('VtuberMember.VtuberName_JP',Name)
    .orWhere('VtuberGroup.VtuberGroupName',Name)
  })
  .orderBy('UploadDate',"desc")
  .limit(Limit)

  data.forEach(i => {
    delete i.id
    delete i.VtuberGroup_id
    delete i.VtuberMember_id
    i.Viewers = parseInt(i.Viewers,10)
  });

  result(null,data)
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
    this.where('VtuberMember.VtuberName', Name)
    .orWhere('VtuberMember.VtuberName_EN',Name)
    .orWhere('VtuberMember.VtuberName_JP',Name)
    .orWhere('VtuberGroup.VtuberGroupName',Name)
  })
  data.forEach(i => {
    delete i.id
    delete i.VtuberGroup_id
    delete i.VtuberMember_id
  });

  result(null,data)
  } catch (error) {
    console.log(error)
    result({kind:"Error kntl"},null)
  }
};


module.exports = {
  GetMemberAll: GetMemberAll,
  GetMemberName: GetMemberName,
  GetGroupName : GetGroupName,
  GetGroupAll: GetGroupAll,
  GetYtLivestream : GetYtLivestream,
  GetTwitter : GetTwitter,
  GetTBilibili : GetTBilibili,
  GetLiveBilibili: GetLiveBilibili,
  GetSpaceBiliBIli: GetSpaceBiliBIli,
  Getsubscriber: Getsubscriber,
}