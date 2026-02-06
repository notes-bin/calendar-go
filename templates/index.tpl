<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>日历订阅中心</title>
    <style>
        * { margin:0; padding:0; box-sizing:border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; }
        body { 
            max-width: 700px; 
            margin: 2rem auto; 
            padding: 0 1rem; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container {
            background: #fff;
            border-radius: 20px;
            padding: 2rem;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
        }
        h2 { 
            color: #333; 
            font-size: 1.8rem; 
            margin-bottom: 0.5rem; 
            text-align: center;
        }
        .subtitle {
            text-align: center;
            color: #666;
            margin-bottom: 1.5rem;
            font-size: 0.9rem;
        }
        .list-container {
            border: 1px solid #e8e8e8;
            border-radius: 12px;
            overflow: hidden;
        }
        .list-item {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 1.2rem 1.5rem;
            border-bottom: 1px solid #e8e8e8;
            transition: all 0.3s ease;
            gap: 1rem;
        }
        .list-item:last-child {
            border-bottom: none;
        }
        .list-item:hover {
            background: #f8f9ff;
        }
        .item-info {
            flex: 1;
            min-width: 0;
        }
        .item-name {
            color: #333;
            font-size: 1.05rem;
            font-weight: 600;
            margin-bottom: 0.3rem;
            word-break: break-word;
        }
        .item-desc {
            color: #999;
            font-size: 0.85rem;
            line-height: 1.4;
            word-break: break-word;
        }
        .btn { 
            display: inline-flex; 
            align-items: center; 
            justify-content: center;
            padding: 0.55rem 1.3rem; 
            border-radius: 8px; 
            text-decoration: none; 
            font-weight: 500;
            font-size: 0.9rem;
            transition: all 0.3s ease;
            white-space: nowrap;
            flex-shrink: 0;
        }
        .btn-sub { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
            color: #fff; 
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
        }
        .btn-sub:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
        }
        @media (max-width: 480px) {
            body {
                margin: 1rem auto;
                padding: 0 0.8rem;
            }
            .container {
                padding: 1.5rem 1rem;
                border-radius: 16px;
            }
            h2 {
                font-size: 1.5rem;
            }
            .subtitle {
                font-size: 0.85rem;
                margin-bottom: 1.2rem;
            }
            .list-item {
                padding: 1rem;
                flex-direction: column;
                align-items: flex-start;
                gap: 0.8rem;
            }
            .item-info {
                width: 100%;
            }
            .item-name {
                font-size: 1rem;
            }
            .item-desc {
                font-size: 0.8rem;
            }
            .btn {
                width: 100%;
                padding: 0.7rem;
                font-size: 0.95rem;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>📅 日历订阅中心</h2>
        <p class="subtitle">点击订阅自动同步到 iPhone / Mac 日历</p>

        <div class="list-container">
        {{ range $key, $sub := .subs }}
        <div class="list-item">
            <div class="item-info">
                <div class="item-name">{{ $sub.Name }}</div>
                <div class="item-desc">{{ $sub.Desc }}</div>
            </div>
            <a href="/subscribe/{{ $key }}" class="btn btn-sub">订阅</a>
        </div>
        {{ end }}
        </div>
    </div>
</body>
</html>