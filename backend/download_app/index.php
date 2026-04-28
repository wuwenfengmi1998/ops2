<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>OPS 运营管理系统</title>
  <style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }
    body {
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    }
    .container {
      background: #fff;
      border-radius: 24px;
      padding: 60px 40px;
      text-align: center;
      box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
      max-width: 400px;
      width: 90%;
    }
    .logo {
      width: 120px;
      height: 120px;
      margin-bottom: 30px;
      object-fit: contain;
    }
    .title {
      font-size: 28px;
      font-weight: 600;
      color: #333;
      margin-bottom: 10px;
    }
    .subtitle {
      font-size: 14px;
      color: #999;
      margin-bottom: 40px;
    }
    .download-btn {
      display: inline-block;
      padding: 16px 48px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: #fff;
      font-size: 18px;
      font-weight: 500;
      border-radius: 50px;
      text-decoration: none;
      transition: transform 0.2s, box-shadow 0.2s;
      box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
    }
    .download-btn:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
    }
    .download-btn:active {
      transform: translateY(0);
    }
    .loading {
      color: #666;
      font-size: 14px;
    }
    .error {
      color: #ff4d4f;
      font-size: 14px;
      margin-top: 20px;
    }
    .version {
      margin-top: 20px;
      font-size: 12px;
      color: #999;
    }
  </style>
</head>
<body>
  <div class="container">
    <img src="logo.png" alt="Logo" class="logo">
    <h1 class="title">OPS 运营管理系统</h1>
    <p class="subtitle">点击下方按钮下载移动端 APP</p>
    
    <div id="content">
      <a href="javascript:void(0)" class="download-btn" onclick="downloadLatest()">
        立即下载
      </a>
      <p class="version" id="version"></p>
    </div>
    
    <div id="loading" style="display:none;">
      <p class="loading">正在获取最新版本...</p>
    </div>
  </div>

  <script>
    function downloadLatest() {
      const content = document.getElementById('content');
      const loading = document.getElementById('loading');
      const version = document.getElementById('version');
      
      content.style.display = 'none';
      loading.style.display = 'block';
      
      fetch('download.php')
        .then(response => response.json())
        .then(data => {
          loading.style.display = 'none';
          if (data.success && data.file) {
            version.textContent = data.filename + ' (' + data.date + ')';
            window.location.href = data.file;
          } else {
            content.style.display = 'block';
            alert(data.message || '未找到 APK 文件');
          }
        })
        .catch(err => {
          loading.style.display = 'none';
          content.style.display = 'block';
          alert('下载失败: ' + err.message);
        });
    }
  </script>
</body>
</html>
