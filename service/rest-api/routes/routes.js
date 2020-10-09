module.exports = app => {
    const Controller = require("../controllers/controller.js");

    //NoG stand for VtuberName OR VtuberGroup
    //youtube,twitter,bilibili,space.bilibili and subscriber have limit parameter,by default limit was 30

    // Select all vtuber data
    app.get("/member", Controller.memberAll); 
    //example /member (string)

    // Select a single vtuber with vtubername[VtuberName/VtuberName_EN/VtuberName_JP]
    app.get("/member/:name", Controller.memberName); 
    //example /member/kanochi (array,separated with commas)

    // Select all vtubergroup
    app.get("/group", Controller.groupAll);
    //example /group (string)
    
    // Retrieve a single vtuber with vtubername[VtuberName/VtuberName_EN/VtuberName_JP]
    app.get("/group/:group", Controller.groupName);
    //example /group/hololive,hanayori (array,separated with commas)

    // Select all Yt
    app.get("/youtube/:nog/:status", Controller.ytlivestream);
    //example /youtube/hololive/past or /youtube/kanochi/upcoming

    // Select vtuber fanart
    app.get("/twitter/:nog", Controller.twitterd)
    //example twitter/hololive 

    //T.bilibili
    app.get("tbilibili/:nog", Controller.tBilibili)
    //example tbilibili/hololive

    //live.bilibili
    app.get("livebilibili/:nog/:status", Controller.liveBilibili)
    //example livebilibili/nijisanji/past

    //space.bilibili
    app.get("spacebilibili/:nog", Controller.spaceBilibili)
    //example spacebilibili/hololive/ or spacebilibili/fbk/

     //subscriber Vtuber count
    app.get("subscriber/:nog", Controller.subscriber)
    //example subscriber/nijisanji or subscriber/fbk

  };