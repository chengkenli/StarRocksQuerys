// 1. 定义一个全局变量 API_APPID（可以是数组、对象或变量）
let API_APPID = "";
// 2. 获取 select 元素并监听 change 事件
const selectBox = document.getElementById("dynamicSelect");
// 初始化时设置默认值（如果是空就设为 'sr-adhoc'）
API_APPID = selectBox.value || "sr-adhoc";
selectBox.addEventListener("change", (event) => {
  // 3. 更新 API_APPID 的值（可以是多个元素或变量）
  API_APPID = event.target.value;
  console.log("Current API_APPID:", API_APPID); // 打印到控制台
  document.querySelectorAll("[data-API_APPID]").forEach((el) => {
    el.textContent = API_APPID; // 或 el.value = API_APPID;
  });
});

function getAppID() {
  return API_APPID;
}
