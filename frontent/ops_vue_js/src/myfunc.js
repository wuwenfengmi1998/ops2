
export const myfuncs = {

	themeStorageKey:"tablerTheme",
	defaultTheme:"light",

  test(){
    console.log("myfuncs test ok");
  },

  	//临时保存的数据，浏览器专属
  	saveT(key,data){
		sessionStorage.setItem(key, data)
	},
	loadT(key){
		return sessionStorage.getItem(key)
	},
	deleT(key){
		sessionStorage.removeItem(key)
	},
	saveJsonT(key,data){
		this.saveT(key,JSON.stringify(data))
	},

	loadJsonT(key){
		const js_data=this.loadT(key)
		if(js_data){
			return JSON.parse(js_data)
		}else{
			return null
		}

	},


	save(key,data){
		localStorage.setItem(key, data)
	},
	load(key){
		return localStorage.getItem(key)
	},
	dele(key){
		localStorage.removeItem(key)
	},

	saveJson(key,data){
		this.save(key,JSON.stringify(data))
	},

	loadJson(key){
		const js_data=this.load(key)
		if(js_data){
			return JSON.parse(js_data)
		}else{
			return null
		}

	},

	getThemefromStorge() {
		const storedTheme = this.load(this.themeStorageKey);
		return storedTheme ? storedTheme : this.defaultTheme;
	},

	setTheme(selectedTheme,save) {
		if(save){
			this.save(this.themeStorageKey, selectedTheme); // 保存到本地存储
		}
		if (selectedTheme === 'dark') {
			document.body.setAttribute("data-bs-theme", selectedTheme); // 暗色模式
		} else {
			document.body.removeAttribute("data-bs-theme"); // 亮色模式（移除属性）
		}
	},

	isValidEmail(email) {
		// 定义邮箱的正则表达式
		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		// 使用正则表达式测试邮箱
		return emailRegex.test(email);
	}
}
