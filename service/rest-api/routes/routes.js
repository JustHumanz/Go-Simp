module.exports = app => {
    const Controller = require("../controllers/controller.js");

    // Select all vtuber data
    app.get("/member", Controller.memberAll);

    // Select a single vtuber with vtubername[VtuberName/VtuberName_EN/VtuberName_JP/VtuberName_CN]
    app.get("/member/:name", Controller.memberName);

    // Select all vtubergroup
    app.get("/group", Controller.groupAll);
    
    // Retrieve a single vtuber with vtubername[VtuberName/VtuberName_EN/VtuberName_JP/VtuberName_CN]
    app.get("/group/:group", Controller.groupName);

    // Select all Yt *nog(name or group)
    app.get("/:nog/youtube/:status", Controller.ytlivestream);

    // Select vtuber fanart
    app.get("/:nog/twitter", Controller.twitterd)

    //T.bilibili
    app.get("/:nog/tbilibili", Controller.tBilibili)

    //live.bilibili
    app.get("/:nog/livebilibili/:status", Controller.liveBilibili)

    //space.bilibili
    app.get("/:nog/spacebilibili", Controller.spaceBilibili)

     //subscriber Vtuber count
    app.get("/:nog/subscriber", Controller.subscriber)

  };