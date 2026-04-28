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

// 按文件修改时间排序，取最新的
$latestFile = null;
$latestMtime = 0;

foreach ($apkFiles as $file) {
    $mtime = filemtime($file);
    if ($mtime > $latestMtime) {
        $latestMtime = $mtime;
        $latestFile = $file;
    }
}

if ($latestFile && file_exists($latestFile)) {
    $filename = basename($latestFile);
    $dateStr = date('Y-m-d H:i:s', $latestMtime);

    // 返回相对于 download_app 目录的路径（需要上跳一级）
    echo json_encode([
        'success' => true,
        'file' => '../' . rawurlencode($filename),
        'filename' => $filename,
        'date' => $dateStr
    ]);
} else {
    echo json_encode([
        'success' => false,
        'message' => 'APK 文件不存在'
    ]);
}
