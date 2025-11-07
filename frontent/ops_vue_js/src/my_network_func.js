import axios from "axios";
import { myfuncs } from "./myfunc";

export const my_network_func = {
  post_json(path, json, callback) {
    var head_path = "/api/v1";
    //把cookie插入json
    var data = {};
    data["data"] = json;
    var cookie = myfuncs.load_json("cookie");
    if (cookie) {
      data["cookie"] = cookie.Value;
    }
    var re_data = {};

    axios
      .post(head_path + path, data, {
        headers: {
          "Content-Type": "application/json",
        },
      })
      .then((response) => {
        //console.log(response)
        re_data["statusCode"] = response.status;
        //载入服务器返回的数据
        if (response.data) {
          re_data["data"] = response.data;
          //自动保存服务器发送的cookie
          if (response.data.cookie) {
            if (response.data.cookie.Value == "") {
              myfuncs.dele("cookie");
            } else {
              myfuncs.save_json("cookie", response.data.cookie);
            }
          }
        }

        callback(re_data);
      })
      .catch((error) => {
        re_data["statusCode"] = -1;
        re_data["error"] = error;
        callback(re_data);
      });
  },
};
