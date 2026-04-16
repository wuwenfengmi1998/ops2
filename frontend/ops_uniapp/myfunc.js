export const myfunc = {

	themeStorageKey:"tablerTheme",
	defaultTheme:"light",
	
	save(key,data){
		uni.setStorageSync(key, data)
	},
	load(key){
		return uni.getStorageSync(key)
	},
	dele(key){
		uni.removeStorageSync(key)
	},
	
	save_json(key,data){
		this.save(key,JSON.stringify(data))
	},
	
	load_json(key){
		var js_data=this.load(key)
		if(js_data){
			return JSON.parse(js_data)
		}else{
			return null
		}
		
	},

	getThemefromStorge() {
		var storedTheme = this.load(this.themeStorageKey);
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


	formatDate(date) {
		var year = date.getFullYear();
		var month = (date.getMonth() + 1).toString().padStart(2, '0');
		var day = date.getDate().toString().padStart(2, '0');
		return `${year}-${month}-${day}`;
	},
	test(tttt) {
		console.log(tttt)
	},
	isValidEmail(email) {
		// 定义邮箱的正则表达式 
		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		// 使用正则表达式测试邮箱 
		return emailRegex.test(email);
	}
}