const model = require("../models/model.js");

exports.MemberAll = (req,res) => {
    model.GetMemberAll(req.query.region,(err, data) => {
        if (err) {
            if (err.kind === "not_found") {
                res.status(404).send({
                message: `Not found VtuberName with name ${req.params.name} and region ${req.query.region}`
                });
            } else {
                res.status(500).send({
                message: "Error retrieving VtuberName with name " + req.params.name
                });
            }
        } else {
            res.send(data);
        }
    });
};

exports.MemberName = (req, res) => {
    let region = null
    if (req.query.region != null){
        region = req.query.region.split(",")
    } else {
        region = null
    }
    model.GetMemberName(req.params.nog.split(","),region, (err, data) => {
    if (err) {
    if (err.kind === "not_found") {
        res.status(404).send({
        message: `Not found VtuberName with name ${req.params.nog}.`
        });
    } else {
        res.status(500).send({
        message: "Error retrieving VtuberName with name " + req.params.nog
        });
    }
    } else res.send(data);
});
};
  

exports.GroupAll = (_,res) => {
    model.GetGroupAll((err, data) => {
        if (err) {
            if (err.kind === "not_found") {
                res.status(404).send({
                message: `Not found Vtuber Group with name ${req.params.name}`
                });
            } else {
                res.status(500).send({
                message: "Error retrieving Vtuber Group with name " + req.params.name
                });
            }
        } else {
            res.send(data);
        }
    });
};
    
exports.GroupName = (req, res) => {
    model.GetGroupName(req.params.group.split(","), (err, data) => {
        if (err) {
        if (err.kind === "not_found") {
            res.status(404).send({
            message: `Not found GroupName with name ${req.params.group}.`
            });
        } else {
            res.status(500).send({
            message: "Error retrieving GroupName with name " + req.params.group
            });
        }
        } else res.send(data);
    });
};


exports.YtliveStream = (req, res) => {
    const Limit =  req.query.limit || 30
    if (Limit >= 100 ){
        res.status(401).send({
            message: `out of limit`
        });
        return
    } 
    model.GetYtLivestream(req.params.nog.split(","),req.params.status, Limit,(err, data) => {
        if (err) {
        if (err.kind === "not_found") {
            res.status(404).send({
            message: `It looks like ${req.params.nog} doesn't have a ${req.params.status} stream right now .`
            });
        } else {
            res.status(500).send({
            message: "Error retrieving GetYtLivestream with name " + req.params.nog
            });
        }
        } else res.send(data);
    });
};

exports.Twitterd = (req, res) => {
    const Limit =  req.query.limit || 30
    if (Limit >= 300 ){
        res.status(401).send({
            message: `out of limit`
        });
        return
    }
    model.GetTwitter(req.params.nog.split(","), Limit,(err, data) => {
        if (err) {
            if (err.kind === "not_found") {
                res.status(404).send({
                message: `Not found GetTwitter with name ${req.params.nog}.`
                });
            } else {
                res.status(500).send({
                message: "Error retrieving GetTwitter with name " + req.params.nog
                });
            }
        } else {
            res.send(data);
        } 
    });
};

exports.Tbilibili = (req, res) => {
    const Limit =  req.query.limit || 30
    if (Limit >= 300 ){
        res.status(401).send({
            message: `out of limit`
        });
        return
    }
    model.GetTBilibili(req.params.nog.split(","), Limit,(err, data) => {
        if (err) {
        if (err.kind === "not_found") {
            res.status(404).send({
            message: `Not found GetTBilibili with name ${req.params.nog}.`
            });
        } else {
            res.status(500).send({
            message: "GetTBilibili Error LMAO " + req.params.nog
            });
        }
        } else res.send(data);
    });
};


exports.LiveBilibili = (req, res) => {
    const Limit =  req.query.limit || 10
    if (Limit >= 30 ){
        res.status(401).send({
            message: `out of limit`
        });
        return
    }
    model.GetLiveBilibili(req.params.nog.split(","),req.params.status, Limit,(err, data) => {
        if (err) {
        if (err.kind === "not_found") {
            res.status(404).send({
            message: `It looks like ${req.params.nog} doesn't have a ${req.params.status} stream right now .`
            });
        } else {
            res.status(500).send({
            message: "Error retrieving GetYtLivestream with name " + req.params.nog
            });
        }
        } else res.send(data);
    });
};

exports.SpaceBilibili = (req, res) => {
    const Limit =  req.query.limit || 30
    if (Limit >= 60 ){
        res.status(401).send({
            message: `out of limit`
        });
        return
    }
    model.GetSpaceBiliBili(req.params.nog.split(","), Limit,(err, data) => {
        if (err) {
        if (err.kind === "not_found") {
            res.status(404).send({
            message: `It looks like ${req.params.nog} doesn't have a ${req.params.status} stream right now .`
            });
        } else {
            res.status(500).send({
            message: "Error retrieving GetYtLivestream with name " + req.params.nog
            });
        }
        } else res.send(data);
    });
};


exports.Subscriber = (req, res) => {
    model.Getsubscriber(req.params.nog.split(","),(err, data) => {
        if (err) {
        if (err.kind === "not_found") {
            res.status(404).send({
            message: `It looks like ${req.params.nog} doesn't have a ${req.params.status} stream right now .`
            });
        } else {
            res.status(500).send({
            message: "Error retrieving GetYtLivestream with name " + req.params.nog
            });
        }
        } else res.send(data);
    });
};


exports.Channel = (req, res) => {
    model.GetDiscordChannel(req.params.id, (err, data) => {
        if (err) {
            if (err.kind === "not_found") {
                res.status(404).send({
                message: `Not found GroupName with name ${req.params.id}.`
                });
            } else {
                res.status(500).send({
                message: "Error retrieving GroupName with name " + req.params.id
                });
            }
        } else {
            res.send(data);
        } 
    });
};