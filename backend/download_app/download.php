<?php
header('Content-Type: application/json; charset=utf-8');

// 扫描当前目录下的所有 APK 文件
$apkFiles = glob('*.apk');

if (empty($apkFiles)) {
    echo json_encode([
        'success' => false,
        'message' => '未找到 APK 文件'
    ]);
    exit;
}

// 解析 APK 文件名中的日期，格式: __UNI__8A0DE5E__20260424215857.apk
// 日期部分: 20260424215857 (YYYYMMDDHHmmss)
$latestFile = null;
$latestDate = 0;

foreach ($apkFiles as $file) {
    $filename = basename($file);
    
    // 提取日期部分 (文件名最后一段，移除 .apk 后缀)
    if (preg_match('/(\d{14})\.apk$/i', $filename, $matches)) {
        $date = intval($matches[1]);
        if ($date > $latestDate) {
            $latestDate = $date;
            $latestFile = $file;
        }
    }
}

if ($latestFile && file_exists($latestFile)) {
    // 返回当前目录下的 APK 文件路径
    $relativePath = './' . basename($latestFile);
    
    echo json_encode([
        'success' => true,
        'file' => $relativePath,
        'date' => substr($latestDate, 0, 8) . ' ' . substr($latestDate, 8, 2) . ':' . substr($latestDate, 10, 2) . ':' . substr($latestDate, 12, 2)
    ]);
} else {
    echo json_encode([
        'success' => false,
        'message' => 'APK 文件不存在'
    ]);
}
