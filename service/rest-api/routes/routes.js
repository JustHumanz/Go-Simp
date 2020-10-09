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
    
    // Retrieve a single vtuber with vtubername[VtuberName/VtuberName_EN/VtuberName_JP/VtuberName_CN]
    app.get("/group/:group", Controller.groupName);
    //example /group/hololive,hanayori (array,separated with commas)

    // Select all Yt
    app.get("/:nog/youtube/:status", Controller.ytlivestream);
    //example /hololive/youtube/past or /kanochi/youtube/upcoming

    // Select vtuber fanart
    app.get("/:nog/twitter", Controller.twitterd)
    //example /hololive/twitter 

    //T.bilibili
    app.get("/:nog/tbilibili", Controller.tBilibili)
    //example /hololive/tbilibili 

    //live.bilibili
    app.get("/:nog/livebilibili/:status", Controller.liveBilibili)
    //example /hololive/livebilibili/past

    //space.bilibili
    app.get("/:nog/spacebilibili", Controller.spaceBilibili)
    //example /hololive/spacebilibili or /fbk/spacebilibili

     //subscriber Vtuber count
    app.get("/:nog/subscriber", Controller.subscriber)
    //example /hololive/subscriber or /fbk/subscriber

  };