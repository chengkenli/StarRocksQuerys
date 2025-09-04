window.setupAutoRefresh = function (interval = 5000) {
  if (window.refreshInterval) {
    clearInterval(window.refreshInterval);
  }
  window.refreshInterval = setInterval(window.fetchData, interval);
  window.refreshInterval = setInterval(window.fetchData_submit, interval);
  window.refreshInterval = setInterval(window.fetchData_brokerload, 10000);
  window.refreshInterval = setInterval(window.fetchData_resourcegroup, 10000);
};
