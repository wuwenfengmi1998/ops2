import { myfunc } from "./myfunc";

export const my_network_func = {
	host: "",
	port: 0,
	head_path: "/api/v1",

	post_json(path, json,callback) {
		//把cookie插入json
		var data={}
		data['data']=json
		var cookie =myfunc.load_json("cookie")
		if(cookie)
		{
			data['cookie']=cookie
		}
		var re_data = {}
		uni.request({
			header: {
				'Content-Type': 'application/json'
			},
			url: this.head_path + path,
			method: 'POST',
			data: data,
			timeout: 10000,
			success(res) {

				re_data["statusCode"] = res.statusCode
				//载入服务器返回的数据
				if(res.data){
					re_data["data"] = res.data
				}
				//自动保存服务器发送的cookie
				if(res.cookie){
					if(res.cookie.Value=="")
					{
						myfunc.dele("cookie")
					}else{
						myfunc.save_json("cookie", res.cookie)
					}
					
				}
				callback(re_data)
				
			},
			fail() {
				re_data["statusCode"] = -1
				callback(re_data)
				
			}

		});
	},
}