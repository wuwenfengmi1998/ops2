
export const mytablercon = {
  set_bs_theme(color: string){

    switch(color){
      case "dark":
        document.body.setAttribute("data-bs-theme", "dark"); // 暗色模式
        break;
      case "light":
        document.body.setAttribute("data-bs-theme", "light"); // 亮色模式
        break;
      default:
        document.body.removeAttribute("data-bs-theme"); // 跟随系统
    }


  },

}
