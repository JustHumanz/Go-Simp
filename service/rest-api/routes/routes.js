module.exports = app => {
    const Controller = require("../controllers/controller.js");

    //NoG stand for VtuberName OR VtuberGroup
    //youtube,twitter,bilibili,space.bilibili and subscriber have limit parameter,by default limit was 30

    // Select all vtuber data
    app.get("/member", Controller.MemberAll); 
    //example /member (string)

    // Select a single vtuber with vtubername[VtuberName/VtuberName_EN/VtuberName_JP]
    app.get("/member/:nog", Controller.MemberName); 
    //example /member/kanochi (array,separated with commas)

    // Select all vtubergroup
    app.get("/group", Controller.GroupAll);
    //example /group (string)
    
    // Retrieve a single vtuber with vtubername[VtuberName/VtuberName_EN/VtuberName_JP]
    app.get("/group/:group", Controller.GroupName);
    //example /group/hololive,hanayori (array,separated with commas)

    // Select all Yt
    app.get("/youtube/:nog/:status", Controller.YtliveStream);
    //example /youtube/hololive/past or /youtube/kanochi/upcoming

    // Select vtuber fanart
    app.get("/twitter/:nog", Controller.Twitterd)
    //example twitter/hololive 

    //T.bilibili
    app.get("/tbilibili/:nog", Controller.Tbilibili)
    //example tbilibili/hololive

    //live.bilibili
    app.get("/livebilibili/:nog/:status", Controller.LiveBilibili)
    //example livebilibili/nijisanji/past

    //space.bilibili
    app.get("/spacebilibili/:nog", Controller.SpaceBilibili)
    //example spacebilibili/hololive/ or spacebilibili/fbk/

     //subscriber Vtuber count
    app.get("/subscriber/:nog", Controller.Subscriber)
    //example subscriber/nijisanji or subscriber/fbk

    app.get("/channel/:id",Controller.Channel)

    app.get('/doc',function(req,res) {
      res.sendFile(__dirname +'/doc.html');
    });
  };