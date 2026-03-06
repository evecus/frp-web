package main

const uiHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>frpc-web</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
:root{
  --primary:#6366f1;--primary-light:#818cf8;--primary-dark:#4f46e5;
  --success:#10b981;--danger:#ef4444;--warning:#f59e0b;--info:#3b82f6;
  --bg:#f1f5f9;--sidebar:#ffffff;--card:#ffffff;
  --text:#1e293b;--text-muted:#64748b;--border:#e2e8f0;
  --shadow:0 1px 3px rgba(0,0,0,.08),0 1px 2px rgba(0,0,0,.04);
  --shadow-md:0 4px 6px -1px rgba(0,0,0,.07),0 2px 4px -1px rgba(0,0,0,.04);
  --radius:12px;--radius-sm:8px;
}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:var(--bg);color:var(--text);min-height:100vh}

/* ── Layout ── */
#login-page{display:flex;align-items:center;justify-content:center;min-height:100vh;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%)}
#app{display:flex;min-height:100vh}
.sidebar{width:240px;background:var(--sidebar);border-right:1px solid var(--border);display:flex;flex-direction:column;position:fixed;top:0;left:0;height:100vh;z-index:100}
.main{margin-left:240px;flex:1;display:flex;flex-direction:column;min-height:100vh}
.topbar{height:60px;background:#fff;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;padding:0 24px;position:sticky;top:0;z-index:50}
.content{padding:24px;flex:1}

/* ── Login ── */
.login-card{background:#fff;border-radius:20px;padding:40px;width:380px;box-shadow:0 20px 60px rgba(0,0,0,.2)}
.login-logo{display:flex;align-items:center;gap:12px;margin-bottom:32px;justify-content:center}
.login-logo .logo-box{width:44px;height:44px;background:linear-gradient(135deg,var(--primary),var(--primary-dark));border-radius:10px;display:flex;align-items:center;justify-content:center;color:#fff;font-weight:800;font-size:13px}
.login-logo h1{font-size:22px;font-weight:700;color:var(--text)}
.login-card h2{font-size:16px;color:var(--text-muted);text-align:center;margin-bottom:28px;font-weight:400}
.field{margin-bottom:16px}
.field label{display:block;font-size:13px;font-weight:500;color:var(--text-muted);margin-bottom:6px}
.field input{width:100%;padding:10px 14px;border:1.5px solid var(--border);border-radius:var(--radius-sm);font-size:14px;transition:.2s;outline:none;color:var(--text)}
.field input:focus{border-color:var(--primary);box-shadow:0 0 0 3px rgba(99,102,241,.1)}
.login-hint{font-size:12px;color:var(--text-muted);text-align:center;margin-top:16px}

/* ── Sidebar ── */
.sidebar-logo{padding:20px 16px;display:flex;align-items:center;gap:10px;border-bottom:1px solid var(--border)}
.sidebar-logo .logo-box{width:36px;height:36px;background:linear-gradient(135deg,var(--primary),var(--primary-dark));border-radius:8px;display:flex;align-items:center;justify-content:center;color:#fff;font-weight:800;font-size:11px;flex-shrink:0}
.sidebar-logo span{font-weight:700;font-size:16px;color:var(--text)}
.sidebar-logo small{font-size:11px;color:var(--text-muted);display:block}
.sidebar-nav{flex:1;padding:12px 8px;overflow-y:auto}
.nav-section{font-size:10px;font-weight:600;color:var(--text-muted);letter-spacing:.08em;text-transform:uppercase;padding:8px 8px 4px}
.nav-item{display:flex;align-items:center;gap:10px;padding:9px 12px;border-radius:var(--radius-sm);cursor:pointer;color:var(--text-muted);font-size:14px;font-weight:500;transition:.15s;margin-bottom:2px;text-decoration:none}
.nav-item:hover{background:#f8fafc;color:var(--text)}
.nav-item.active{background:linear-gradient(135deg,rgba(99,102,241,.12),rgba(99,102,241,.06));color:var(--primary);font-weight:600}
.nav-item svg{flex-shrink:0;opacity:.7}
.nav-item.active svg{opacity:1}
.sidebar-footer{padding:12px 8px;border-top:1px solid var(--border)}

/* ── Topbar ── */
.topbar-left{display:flex;align-items:center;gap:16px}
.page-title{font-size:18px;font-weight:600}
.page-sub{font-size:12px;color:var(--text-muted)}
.topbar-right{display:flex;align-items:center;gap:12px}
.status-badge{display:flex;align-items:center;gap:6px;padding:6px 12px;border-radius:20px;font-size:13px;font-weight:500}
.status-badge.running{background:#dcfce7;color:#16a34a}
.status-badge.stopped{background:#fef2f2;color:#dc2626}
.status-badge.connecting{background:#fef9c3;color:#ca8a04}
.status-dot{width:7px;height:7px;border-radius:50%}
.running .status-dot{background:#16a34a;animation:pulse 2s infinite}
.stopped .status-dot{background:#dc2626}
.connecting .status-dot{background:#ca8a04;animation:pulse 1s infinite}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:.4}}
.ctrl-btns{display:flex;gap:6px}

/* ── Buttons ── */
.btn{display:inline-flex;align-items:center;gap:6px;padding:8px 16px;border-radius:var(--radius-sm);font-size:13px;font-weight:500;cursor:pointer;border:none;transition:.15s;white-space:nowrap}
.btn:disabled{opacity:.5;cursor:not-allowed}
.btn-primary{background:linear-gradient(135deg,var(--primary),var(--primary-dark));color:#fff;box-shadow:0 2px 8px rgba(99,102,241,.3)}
.btn-primary:hover:not(:disabled){box-shadow:0 4px 12px rgba(99,102,241,.4);transform:translateY(-1px)}
.btn-success{background:linear-gradient(135deg,#10b981,#059669);color:#fff;box-shadow:0 2px 8px rgba(16,185,129,.3)}
.btn-danger{background:linear-gradient(135deg,#ef4444,#dc2626);color:#fff;box-shadow:0 2px 8px rgba(239,68,68,.3)}
.btn-warning{background:linear-gradient(135deg,#f59e0b,#d97706);color:#fff}
.btn-ghost{background:transparent;color:var(--text-muted);border:1.5px solid var(--border)}
.btn-ghost:hover{background:#f8fafc;color:var(--text)}
.btn-sm{padding:5px 10px;font-size:12px}
.btn-icon{padding:7px;border-radius:var(--radius-sm);background:transparent;color:var(--text-muted);border:1.5px solid var(--border);cursor:pointer;display:flex;align-items:center;transition:.15s}
.btn-icon:hover{background:#f8fafc;color:var(--text)}
.btn-login{width:100%;padding:12px;background:linear-gradient(135deg,var(--primary),var(--primary-dark));color:#fff;border:none;border-radius:var(--radius-sm);font-size:15px;font-weight:600;cursor:pointer;transition:.2s;margin-top:8px}
.btn-login:hover{transform:translateY(-1px);box-shadow:0 8px 20px rgba(99,102,241,.35)}

/* ── Cards ── */
.card{background:var(--card);border-radius:var(--radius);box-shadow:var(--shadow);border:1px solid var(--border)}
.card-header{padding:16px 20px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between}
.card-title{font-size:14px;font-weight:600;color:var(--text);display:flex;align-items:center;gap:8px}
.card-body{padding:20px}

/* ── Stat cards ── */
.stats-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:16px;margin-bottom:24px}
.stat-card{background:var(--card);border-radius:var(--radius);padding:20px;border:1px solid var(--border);position:relative;overflow:hidden}
.stat-card::before{content:'';position:absolute;top:0;left:0;right:0;height:3px}
.stat-card.blue::before{background:linear-gradient(90deg,#6366f1,#818cf8)}
.stat-card.green::before{background:linear-gradient(90deg,#10b981,#34d399)}
.stat-card.red::before{background:linear-gradient(90deg,#ef4444,#f87171)}
.stat-card.orange::before{background:linear-gradient(90deg,#f59e0b,#fbbf24)}
.stat-label{font-size:12px;color:var(--text-muted);font-weight:500;margin-bottom:8px}
.stat-value{font-size:28px;font-weight:700;color:var(--text)}
.stat-sub{font-size:11px;color:var(--text-muted);margin-top:4px}

/* ── Table ── */
.table-wrap{overflow-x:auto}
table{width:100%;border-collapse:collapse;font-size:13px}
th{text-align:left;padding:10px 14px;font-size:11px;font-weight:600;color:var(--text-muted);text-transform:uppercase;letter-spacing:.05em;border-bottom:2px solid var(--border);background:#fafafa}
td{padding:11px 14px;border-bottom:1px solid var(--border);color:var(--text);vertical-align:middle}
tr:last-child td{border-bottom:none}
tr:hover td{background:#fafafe}
.empty-row td{text-align:center;padding:40px;color:var(--text-muted)}

/* ── Badges ── */
.badge{display:inline-flex;align-items:center;padding:2px 8px;border-radius:20px;font-size:11px;font-weight:600;text-transform:uppercase}
.badge-tcp{background:#ede9fe;color:#7c3aed}
.badge-udp{background:#fef3c7;color:#92400e}
.badge-http{background:#dbeafe;color:#1d4ed8}
.badge-https{background:#dcfce7;color:#166534}
.badge-running{background:#dcfce7;color:#166534}
.badge-stopped{background:#fef2f2;color:#991b1b}
.badge-error{background:#fef2f2;color:#991b1b}
.badge-active{background:linear-gradient(135deg,rgba(99,102,241,.15),rgba(99,102,241,.08));color:var(--primary);border:1px solid rgba(99,102,241,.2)}

/* ── Forms ── */
.form-grid{display:grid;grid-template-columns:1fr 1fr;gap:16px}
.form-grid.cols3{grid-template-columns:1fr 1fr 1fr}
.form-field{display:flex;flex-direction:column;gap:5px}
.form-field.full{grid-column:1/-1}
.form-field label{font-size:12px;font-weight:600;color:var(--text-muted);text-transform:uppercase;letter-spacing:.05em}
.form-field input,.form-field select,.form-field textarea{padding:9px 12px;border:1.5px solid var(--border);border-radius:var(--radius-sm);font-size:13px;color:var(--text);outline:none;transition:.2s;background:#fff}
.form-field input:focus,.form-field select:focus{border-color:var(--primary);box-shadow:0 0 0 3px rgba(99,102,241,.1)}
.form-hint{font-size:11px;color:var(--text-muted);margin-top:2px}

/* ── Profile cards ── */
.profiles-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:16px}
.profile-card{background:var(--card);border:2px solid var(--border);border-radius:var(--radius);padding:20px;cursor:pointer;transition:.2s;position:relative}
.profile-card:hover{border-color:var(--primary-light);box-shadow:var(--shadow-md)}
.profile-card.active{border-color:var(--primary);box-shadow:0 0 0 3px rgba(99,102,241,.1)}
.profile-name{font-size:15px;font-weight:600;margin-bottom:6px;display:flex;align-items:center;gap:8px}
.profile-addr{font-size:12px;color:var(--text-muted);font-family:monospace}
.profile-user{font-size:12px;color:var(--text-muted);margin-top:4px}
.profile-actions{display:flex;gap:6px;margin-top:14px}

/* ── Modal ── */
.modal-bg{position:fixed;inset:0;background:rgba(0,0,0,.4);backdrop-filter:blur(4px);z-index:1000;display:flex;align-items:center;justify-content:center;opacity:0;pointer-events:none;transition:.2s}
.modal-bg.open{opacity:1;pointer-events:all}
.modal{background:#fff;border-radius:16px;width:500px;max-width:95vw;box-shadow:0 25px 50px rgba(0,0,0,.2);transform:scale(.95);transition:.2s}
.modal-bg.open .modal{transform:scale(1)}
.modal-head{padding:20px 24px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between}
.modal-title{font-size:16px;font-weight:600}
.modal-body{padding:24px}
.modal-foot{padding:16px 24px;border-top:1px solid var(--border);display:flex;justify-content:flex-end;gap:10px}
.close-btn{background:none;border:none;cursor:pointer;color:var(--text-muted);padding:4px;border-radius:6px;line-height:0}
.close-btn:hover{background:#f1f5f9;color:var(--text)}

/* ── Toast ── */
#toast{position:fixed;top:20px;right:20px;z-index:9999;display:flex;flex-direction:column;gap:8px;pointer-events:none}
.toast-item{padding:12px 16px;border-radius:10px;font-size:13px;font-weight:500;box-shadow:0 8px 24px rgba(0,0,0,.12);animation:slideIn .3s ease;pointer-events:all;display:flex;align-items:center;gap:8px;min-width:240px}
.toast-ok{background:#fff;border-left:4px solid var(--success);color:var(--text)}
.toast-err{background:#fff;border-left:4px solid var(--danger);color:var(--text)}
.toast-info{background:#fff;border-left:4px solid var(--info);color:var(--text)}
@keyframes slideIn{from{transform:translateX(100%);opacity:0}to{transform:translateX(0);opacity:1}}
@keyframes fadeOut{to{opacity:0;transform:translateX(20px)}}

/* ── Server info ── */
.server-info{display:flex;align-items:center;gap:8px;padding:8px 12px;background:#f8faff;border:1px solid rgba(99,102,241,.2);border-radius:var(--radius-sm);font-size:13px}
.server-info .label{color:var(--text-muted)}
.server-info .value{font-weight:600;color:var(--primary);font-family:monospace}

/* ── Misc ── */
.hidden{display:none!important}
.divider{height:1px;background:var(--border);margin:16px 0}
.text-muted{color:var(--text-muted)}
.font-mono{font-family:monospace}
.flex{display:flex}.items-center{align-items:center}.gap-2{gap:8px}.gap-3{gap:12px}.justify-between{justify-content:space-between}
.mb-4{margin-bottom:16px}.mb-6{margin-bottom:24px}
.page{display:none}.page.active{display:block}
.section-title{font-size:14px;font-weight:600;color:var(--text);margin-bottom:16px;display:flex;align-items:center;gap:8px}
.spinner{width:16px;height:16px;border:2px solid rgba(255,255,255,.3);border-top-color:#fff;border-radius:50%;animation:spin .6s linear infinite;display:inline-block}
@keyframes spin{to{transform:rotate(360deg)}}
.username-tag{display:inline-flex;align-items:center;gap:4px;padding:3px 8px;background:#ede9fe;color:#7c3aed;border-radius:20px;font-size:11px;font-weight:600}
</style>
</head>
<body>

<!-- Toast -->
<div id="toast"></div>

<!-- Login -->
<div id="login-page">
  <div class="login-card">
    <div class="login-logo">
      <div class="logo-box">FRP</div>
      <h1>frpc-web</h1>
    </div>
    <h2>请登录以继续</h2>
    <div class="field"><label>用户名</label><input id="lu" type="text" value="admin" autocomplete="username"></div>
    <div class="field"><label>密码</label><input id="lp" type="password" placeholder="输入密码" autocomplete="current-password"></div>
    <button class="btn-login" onclick="doLogin()">登 录</button>
    <p class="login-hint">首次使用默认账号 admin / admin</p>
  </div>
</div>

<!-- App -->
<div id="app" class="hidden">
  <aside class="sidebar">
    <div class="sidebar-logo">
      <div class="logo-box">FRP</div>
      <div><span>frpc-web</span><small>客户端管理</small></div>
    </div>
    <nav class="sidebar-nav">
      <div class="nav-section">监控</div>
      <a class="nav-item active" data-page="overview" onclick="nav('overview')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
        概览
      </a>
      <a class="nav-item" data-page="proxies" onclick="nav('proxies')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01"/></svg>
        代理列表
      </a>
      <div class="nav-section">配置</div>
      <a class="nav-item" data-page="servers" onclick="nav('servers')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><rect x="2" y="3" width="20" height="5" rx="1"/><rect x="2" y="10" width="20" height="5" rx="1"/><rect x="2" y="17" width="20" height="5" rx="1"/><circle cx="6" cy="5.5" r="1" fill="currentColor"/><circle cx="6" cy="12.5" r="1" fill="currentColor"/><circle cx="6" cy="19.5" r="1" fill="currentColor"/></svg>
        服务器配置
      </a>
      <div class="nav-section">账户</div>
      <a class="nav-item" data-page="account" onclick="nav('account')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/></svg>
        账户安全
      </a>
    </nav>
    <div class="sidebar-footer">
      <a class="nav-item" onclick="doLogout()">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
        退出登录
      </a>
    </div>
  </aside>

  <main class="main">
    <header class="topbar">
      <div class="topbar-left">
        <div>
          <div class="page-title" id="page-title">概览</div>
          <div class="page-sub" id="server-info-bar">-</div>
        </div>
      </div>
      <div class="topbar-right">
        <div id="status-badge" class="status-badge stopped"><span class="status-dot"></span><span id="status-text">已停止</span></div>
        <div class="ctrl-btns">
          <button class="btn btn-success btn-sm" id="btn-start" onclick="doStart()">
            <svg width="12" height="12" fill="currentColor" viewBox="0 0 24 24"><polygon points="5,3 19,12 5,21"/></svg> 启动
          </button>
          <button class="btn btn-danger btn-sm" id="btn-stop" onclick="doStop()">
            <svg width="12" height="12" fill="currentColor" viewBox="0 0 24 24"><rect x="3" y="3" width="18" height="18" rx="2"/></svg> 停止
          </button>
          <button class="btn btn-warning btn-sm" id="btn-restart" onclick="doRestart()">
            <svg width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg> 重启
          </button>
          <button class="btn-icon" onclick="loadStatus()" title="刷新">
            <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
          </button>
        </div>
        <div id="user-tag" class="username-tag" style="display:none">
          <svg width="10" height="10" fill="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/></svg>
          <span id="user-tag-name"></span>
        </div>
      </div>
    </header>

    <div class="content">
      <!-- Overview -->
      <div id="page-overview" class="page active">
        <div class="stats-grid">
          <div class="stat-card blue"><div class="stat-label">代理总数</div><div class="stat-value" id="s-total">0</div></div>
          <div class="stat-card green"><div class="stat-label">运行中</div><div class="stat-value" id="s-running">0</div></div>
          <div class="stat-card red"><div class="stat-label">错误</div><div class="stat-value" id="s-error">0</div></div>
        </div>
        <div class="card">
          <div class="card-header">
            <div class="card-title">
              <svg width="15" height="15" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01"/></svg>
              代理状态
            </div>
          </div>
          <div class="table-wrap">
            <table><thead><tr><th>名称</th><th>类型</th><th>本地</th><th>远程</th><th>状态</th></tr></thead>
            <tbody id="proxy-status-body"><tr class="empty-row"><td colspan="5">暂无代理</td></tr></tbody></table>
          </div>
        </div>
      </div>

      <!-- Proxies -->
      <div id="page-proxies" class="page">
        <div class="flex justify-between items-center mb-4">
          <div class="section-title">
            <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01"/></svg>
            代理规则
          </div>
          <button class="btn btn-primary" onclick="openProxyModal()">
            <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
            添加代理
          </button>
        </div>
        <div class="card">
          <div class="table-wrap">
            <table><thead><tr><th>名称</th><th>类型</th><th>本地地址</th><th>远程端口 / 域名</th><th>状态</th><th>操作</th></tr></thead>
            <tbody id="proxies-body"><tr class="empty-row"><td colspan="6">暂无代理</td></tr></tbody></table>
          </div>
        </div>
      </div>

      <!-- Servers -->
      <div id="page-servers" class="page">
        <div class="flex justify-between items-center mb-4">
          <div class="section-title">
            <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><rect x="2" y="3" width="20" height="5" rx="1"/><rect x="2" y="10" width="20" height="5" rx="1"/><rect x="2" y="17" width="20" height="5" rx="1"/></svg>
            服务器配置
          </div>
          <button class="btn btn-primary" onclick="openServerModal()">
            <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
            添加服务器
          </button>
        </div>
        <p style="font-size:13px;color:var(--text-muted);margin-bottom:16px">点击卡片激活服务器，启动前需先选择一个服务器配置</p>
        <div class="profiles-grid" id="profiles-grid"></div>
      </div>

      <!-- Account -->
      <div id="page-account" class="page">
        <div class="card" style="max-width:440px">
          <div class="card-header"><div class="card-title">修改登录密码</div></div>
          <div class="card-body">
            <div style="display:flex;flex-direction:column;gap:14px">
              <div class="form-field"><label>当前密码</label><input id="pw-cur" type="password" placeholder="输入当前密码"></div>
              <div class="form-grid">
                <div class="form-field"><label>新密码</label><input id="pw-new" type="password" placeholder="至少4位"></div>
                <div class="form-field"><label>确认新密码</label><input id="pw-cfm" type="password" placeholder="再次输入"></div>
              </div>
              <button class="btn btn-primary" onclick="doChangePw()">保存密码</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
</div>

<!-- Proxy Modal -->
<div class="modal-bg" id="proxy-modal">
  <div class="modal">
    <div class="modal-head">
      <span class="modal-title" id="proxy-modal-title">添加代理</span>
      <button class="close-btn" onclick="closeProxyModal()"><svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
    </div>
    <div class="modal-body">
      <div class="form-grid" style="gap:14px">
        <div class="form-field"><label>代理名称</label><input id="px-name" type="text" placeholder="my-ssh"></div>
        <div class="form-field"><label>类型</label>
          <select id="px-type" onchange="onProxyTypeChange()">
            <option value="tcp">TCP</option><option value="udp">UDP</option>
            <option value="http">HTTP</option><option value="https">HTTPS</option>
          </select>
        </div>
        <div class="form-field"><label>本地 IP</label><input id="px-lip" type="text" placeholder="127.0.0.1"></div>
        <div class="form-field"><label>本地端口</label><input id="px-lport" type="number" placeholder="22"></div>
        <div class="form-field" id="px-rport-field"><label>远程端口</label><input id="px-rport" type="number" placeholder="6022"></div>
        <div class="form-field hidden" id="px-domain-field"><label>自定义域名</label><input id="px-domain" type="text" placeholder="example.com"></div>
        <div class="form-field hidden" id="px-subdomain-field"><label>子域名</label><input id="px-subdomain" type="text" placeholder="my-app"></div>
      </div>
      <div style="margin-top:14px">
        <label style="display:flex;align-items:center;gap:8px;font-size:13px;cursor:pointer;color:var(--text-muted)">
          <input type="checkbox" id="px-disabled"> 禁用此代理
        </label>
      </div>
    </div>
    <div class="modal-foot">
      <button class="btn btn-ghost" onclick="closeProxyModal()">取消</button>
      <button class="btn btn-primary" onclick="saveProxy()">保存</button>
    </div>
  </div>
</div>

<!-- Server Modal -->
<div class="modal-bg" id="server-modal">
  <div class="modal">
    <div class="modal-head">
      <span class="modal-title" id="server-modal-title">添加服务器</span>
      <button class="close-btn" onclick="closeServerModal()"><svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
    </div>
    <div class="modal-body">
      <div class="form-grid" style="gap:14px">
        <div class="form-field full"><label>配置名称</label><input id="sv-name" type="text" placeholder="我的服务器"></div>
        <div class="form-field"><label>服务器地址</label><input id="sv-addr" type="text" placeholder="example.com"></div>
        <div class="form-field"><label>端口</label><input id="sv-port" type="number" placeholder="7000"></div>
        <div class="form-field full"><label>认证 Token</label><input id="sv-token" type="password" placeholder="留空则不验证"></div>
        <div class="form-field full"><label>用户名 <span style="font-size:10px;color:var(--text-muted);font-weight:400">（显示在 frps 面板，用于区分连接来源）</span></label>
          <input id="sv-username" type="text" placeholder="例如: home-pc"></div>
      </div>
      <div id="sv-test-result" style="margin-top:12px;font-size:13px;display:none"></div>
    </div>
    <div class="modal-foot">
      <button class="btn btn-ghost" onclick="testServerConn()" id="btn-test-conn">测试连接</button>
      <button class="btn btn-ghost" onclick="closeServerModal()">取消</button>
      <button class="btn btn-primary" onclick="saveServer()">保存</button>
    </div>
  </div>
</div>

<script>
const g = id => document.getElementById(id);
const esc = s => String(s ?? '').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');

// ── API ──────────────────────────────────────────────────────────────────────
const api = {
  async req(method, path, body) {
    const opts = {method, credentials:'include', headers:{'Content-Type':'application/json'}};
    if (body !== undefined) opts.body = JSON.stringify(body);
    const res = await fetch('/api' + path, opts);
    const data = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(data.error || res.statusText);
    return data;
  },
  get: p => api.req('GET', p),
  post: (p, b) => api.req('POST', p, b),
  put: (p, b) => api.req('PUT', p, b),
  del: p => api.req('DELETE', p),
};

// ── Toast ────────────────────────────────────────────────────────────────────
function toast(msg, type='ok') {
  const el = document.createElement('div');
  const icons = {ok:'✓', err:'✕', info:'ℹ'};
  el.className = ` + "`" + `toast-item toast-${type}` + "`" + `;
  el.innerHTML = ` + "`" + `<span>${icons[type]||'•'}</span><span>${esc(msg)}</span>` + "`" + `;
  g('toast').appendChild(el);
  setTimeout(() => { el.style.animation='fadeOut .3s ease forwards'; setTimeout(() => el.remove(), 300); }, 3000);
}

// ── Auth ─────────────────────────────────────────────────────────────────────
async function doLogin() {
  const u = g('lu').value.trim(), p = g('lp').value;
  if (!u || !p) return;
  try {
    await api.post('/login', {Username: u, Password: p});
    g('login-page').classList.add('hidden');
    g('app').classList.remove('hidden');
    loadStatus(); loadProfiles();
    setInterval(loadStatus, 5000);
  } catch(e) { toast(e.message, 'err'); }
}

async function doLogout() {
  await api.post('/logout').catch(()=>{});
  g('app').classList.add('hidden');
  g('login-page').classList.remove('hidden');
  g('lp').value = '';
}

// ── Navigation ───────────────────────────────────────────────────────────────
const pageTitles = {overview:'概览', proxies:'代理列表', servers:'服务器配置', account:'账户安全'};
let curPage = 'overview';

function nav(page) {
  curPage = page;
  document.querySelectorAll('.page').forEach(el => el.classList.remove('active'));
  document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));
  g('page-' + page).classList.add('active');
  document.querySelector(` + "`" + `.nav-item[data-page="${page}"]` + "`" + `).classList.add('active');
  g('page-title').textContent = pageTitles[page] || page;
  if (page === 'proxies') loadProxies();
  if (page === 'servers') loadProfiles();
}

// ── Status ───────────────────────────────────────────────────────────────────
async function loadStatus() {
  try {
    const d = await api.get('/status');
    const badge = g('status-badge');
    const dot = badge.querySelector('.status-dot');
    const txt = g('status-text');
    if (d.running) {
      badge.className = 'status-badge running';
      txt.textContent = '运行中';
    } else {
      badge.className = d.error ? 'status-badge connecting' : 'status-badge stopped';
      txt.textContent = d.error ? '连接失败' : '已停止';
    }
    // Server info bar
    if (d.serverAddr) {
      g('server-info-bar').textContent = ` + "`" + `${d.profileName || ''} — ${d.serverAddr}:${d.serverPort}` + "`" + `;
    }
    // Username tag
    if (d.username) {
      g('user-tag').style.display = 'inline-flex';
      g('user-tag-name').textContent = d.username;
    } else {
      g('user-tag').style.display = 'none';
    }
    // Stats
    const proxies = d.proxies || [];
    g('s-total').textContent = proxies.length;
    g('s-running').textContent = proxies.filter(p => p.status === 'running').length;
    g('s-error').textContent = proxies.filter(p => p.err).length;
    // Table
    const tb = g('proxy-status-body');
    if (!proxies.length) {
      tb.innerHTML = '<tr class="empty-row"><td colspan="5">暂无代理</td></tr>';
    } else {
      tb.innerHTML = proxies.map(p => ` + "`" + `<tr>
        <td class="font-mono">${esc(p.name)}</td>
        <td><span class="badge badge-${p.type||'tcp'}">${(p.type||'tcp').toUpperCase()}</span></td>
        <td class="font-mono">${esc(p.localAddr||'-')}</td>
        <td class="font-mono">${esc(p.remoteAddr||'-')}</td>
        <td><span class="badge ${p.err?'badge-error':p.status==='running'?'badge-running':'badge-stopped'}">${p.err?'错误':p.status==='running'?'运行中':'停止'}</span></td>
      </tr>` + "`" + `).join('');
    }
  } catch(e) {}
}

// ── frpc control ─────────────────────────────────────────────────────────────
async function doStart() {
  try { await api.post('/frpc/start'); toast('启动中...','info'); setTimeout(loadStatus,1000); } catch(e) { toast(e.message,'err'); }
}
async function doStop() {
  try { await api.post('/frpc/stop'); toast('已停止'); setTimeout(loadStatus,500); } catch(e) { toast(e.message,'err'); }
}
async function doRestart() {
  try { await api.post('/frpc/restart'); toast('重启中...','info'); setTimeout(loadStatus,1500); } catch(e) { toast(e.message,'err'); }
}

// ── Proxies ──────────────────────────────────────────────────────────────────
let editingProxy = null;

async function loadProxies() {
  try {
    const list = await api.get('/proxies');
    const tb = g('proxies-body');
    if (!list || !list.length) {
      tb.innerHTML = '<tr class="empty-row"><td colspan="6">暂无代理 — 点击「添加代理」创建</td></tr>'; return;
    }
    tb.innerHTML = list.map(p => {
      const remote = p.remotePort ? ` + "`" + `:${p.remotePort}` + "`" + ` : (p.customDomain || p.subdomain || '-');
      return ` + "`" + `<tr>
        <td class="font-mono">${esc(p.name)}</td>
        <td><span class="badge badge-${p.type||'tcp'}">${(p.type||'tcp').toUpperCase()}</span></td>
        <td class="font-mono">${esc(p.localIP||'127.0.0.1')}:${p.localPort}</td>
        <td class="font-mono">${esc(remote)}</td>
        <td><span class="badge ${p.disabled?'badge-stopped':'badge-running'}">${p.disabled?'禁用':'启用'}</span></td>
        <td><div class="flex gap-2">
          <button class="btn btn-ghost btn-sm" onclick='editProxy(${JSON.stringify(p)})'>编辑</button>
          <button class="btn btn-danger btn-sm" onclick="delProxy('${esc(p.name)}')">删除</button>
        </div></td>
      </tr>` + "`" + `;
    }).join('');
  } catch(e) { toast(e.message,'err'); }
}

function openProxyModal(p) {
  editingProxy = p || null;
  g('proxy-modal-title').textContent = p ? '编辑代理' : '添加代理';
  g('px-name').value = p ? p.name : '';
  g('px-name').disabled = !!p;
  g('px-type').value = p ? (p.type||'tcp') : 'tcp';
  g('px-lip').value = p ? (p.localIP||'127.0.0.1') : '127.0.0.1';
  g('px-lport').value = p ? p.localPort : '';
  g('px-rport').value = p ? (p.remotePort||'') : '';
  g('px-domain').value = p ? (p.customDomain||'') : '';
  g('px-subdomain').value = p ? (p.subdomain||'') : '';
  g('px-disabled').checked = p ? !!p.disabled : false;
  onProxyTypeChange();
  g('proxy-modal').classList.add('open');
}

function editProxy(p) { openProxyModal(p); }

function closeProxyModal() { g('proxy-modal').classList.remove('open'); }

function onProxyTypeChange() {
  const t = g('px-type').value;
  const isHttp = t === 'http' || t === 'https';
  g('px-rport-field').classList.toggle('hidden', isHttp);
  g('px-domain-field').classList.toggle('hidden', !isHttp);
  g('px-subdomain-field').classList.toggle('hidden', !isHttp);
}

async function saveProxy() {
  const name = g('px-name').value.trim();
  const localPort = parseInt(g('px-lport').value);
  if (!name) { toast('请填写代理名称','err'); return; }
  if (!localPort) { toast('请填写本地端口','err'); return; }
  const p = {
    name, type: g('px-type').value,
    localIP: g('px-lip').value.trim() || '127.0.0.1',
    localPort,
    remotePort: parseInt(g('px-rport').value) || 0,
    customDomain: g('px-domain').value.trim(),
    subdomain: g('px-subdomain').value.trim(),
    disabled: g('px-disabled').checked,
  };
  try {
    if (editingProxy) {
      await api.put('/proxies/' + encodeURIComponent(editingProxy.name), p);
      toast('代理已更新');
    } else {
      await api.post('/proxies', p);
      toast('代理已添加');
    }
    closeProxyModal();
    loadProxies();
    loadStatus();
  } catch(e) { toast(e.message,'err'); }
}

async function delProxy(name) {
  if (!confirm(` + "`" + `删除代理 "${name}"?` + "`" + `)) return;
  try { await api.del('/proxies/' + encodeURIComponent(name)); toast('已删除'); loadProxies(); loadStatus(); }
  catch(e) { toast(e.message,'err'); }
}

// ── Profiles (multi-server) ───────────────────────────────────────────────────
let editingProfile = null;
let profilesData = [];

async function loadProfiles() {
  try {
    const d = await api.get('/profiles');
    profilesData = d.profiles || [];
    const activeID = d.activeID;
    const grid = g('profiles-grid');
    if (!profilesData.length) {
      grid.innerHTML = '<p style="color:var(--text-muted);font-size:13px">暂无服务器配置</p>'; return;
    }
    grid.innerHTML = profilesData.map(p => ` + "`" + `
      <div class="profile-card ${p.id === activeID ? 'active' : ''}" onclick="activateProfile('${p.id}')">
        <div class="profile-name">
          ${esc(p.name)}
          ${p.id === activeID ? '<span class="badge badge-active">当前</span>' : ''}
        </div>
        <div class="profile-addr">${esc(p.serverAddr)}:${p.serverPort}</div>
        ${p.username ? ` + "`" + `<div class="profile-user">用户名: <strong>${esc(p.username)}</strong></div>` + "`" + ` : ''}
        <div class="profile-actions">
          <button class="btn btn-ghost btn-sm" onclick="event.stopPropagation();editProfile('${p.id}')">编辑</button>
          <button class="btn btn-ghost btn-sm" onclick="event.stopPropagation();testProfile('${p.id}', this)">测试</button>
          ${profilesData.length > 1 ? ` + "`" + `<button class="btn btn-danger btn-sm" onclick="event.stopPropagation();delProfile('${p.id}')">删除</button>` + "`" + ` : ''}
        </div>
      </div>` + "`" + `).join('');
  } catch(e) { toast(e.message,'err'); }
}

function openServerModal(p) {
  editingProfile = p || null;
  g('server-modal-title').textContent = p ? '编辑服务器' : '添加服务器';
  g('sv-name').value = p ? p.name : '';
  g('sv-addr').value = p ? p.serverAddr : '';
  g('sv-port').value = p ? p.serverPort : 7000;
  g('sv-token').value = '';
  g('sv-username').value = p ? (p.username||'') : '';
  g('sv-test-result').style.display = 'none';
  g('server-modal').classList.add('open');
}

function editProfile(id) {
  const p = profilesData.find(x => x.id === id);
  if (p) openServerModal(p);
}

function closeServerModal() { g('server-modal').classList.remove('open'); }

async function saveServer() {
  const name = g('sv-name').value.trim();
  const serverAddr = g('sv-addr').value.trim();
  if (!name) { toast('请填写配置名称','err'); return; }
  if (!serverAddr) { toast('请填写服务器地址','err'); return; }
  const body = {
    name, serverAddr,
    serverPort: parseInt(g('sv-port').value) || 7000,
    authToken: g('sv-token').value,
    username: g('sv-username').value.trim(),
  };
  try {
    if (editingProfile) {
      await api.put('/profiles/' + editingProfile.id, body);
      toast('已更新');
    } else {
      await api.post('/profiles', body);
      toast('已添加');
    }
    closeServerModal();
    loadProfiles();
  } catch(e) { toast(e.message,'err'); }
}

async function activateProfile(id) {
  try {
    await api.post('/profiles/' + id + '/activate');
    toast('已切换服务器','info');
    loadProfiles();
    loadStatus();
  } catch(e) { toast(e.message,'err'); }
}

async function testProfile(id, btn) {
  const orig = btn.textContent;
  btn.textContent = '测试中...'; btn.disabled = true;
  try {
    const r = await api.post('/profiles/' + id + '/test');
    toast(r.msg, r.ok ? 'ok' : 'err');
  } catch(e) { toast(e.message,'err'); }
  btn.textContent = orig; btn.disabled = false;
}

async function testServerConn() {
  const addr = g('sv-addr').value.trim();
  const port = parseInt(g('sv-port').value) || 7000;
  if (!addr) { toast('请先填写服务器地址','err'); return; }
  // Save temp and test
  const btn = g('btn-test-conn');
  btn.disabled = true; btn.textContent = '测试中...';
  const res = g('sv-test-result');
  res.style.display = 'none';
  try {
    // Create a temporary profile, test it, then delete
    const body = {name:'__tmp__', serverAddr:addr, serverPort:port, authToken: g('sv-token').value};
    const created = await api.post('/profiles', body);
    const r = await api.post('/profiles/' + created.id + '/test');
    await api.del('/profiles/' + created.id);
    res.style.display = 'block';
    res.style.color = r.ok ? 'var(--success)' : 'var(--danger)';
    res.textContent = r.msg;
  } catch(e) {
    res.style.display = 'block';
    res.style.color = 'var(--danger)';
    res.textContent = e.message;
  }
  btn.disabled = false; btn.textContent = '测试连接';
}

async function delProfile(id) {
  if (!confirm('删除此服务器配置?')) return;
  try { await api.del('/profiles/' + id); toast('已删除'); loadProfiles(); }
  catch(e) { toast(e.message,'err'); }
}

// ── Password ──────────────────────────────────────────────────────────────────
async function doChangePw() {
  const cur = g('pw-cur').value, nw = g('pw-new').value, cfm = g('pw-cfm').value;
  if (!cur || !nw) { toast('请填写完整','err'); return; }
  if (nw !== cfm) { toast('两次密码不一致','err'); return; }
  try { await api.post('/password', {Current:cur, New:nw}); toast('密码已修改'); g('pw-cur').value=g('pw-new').value=g('pw-cfm').value=''; }
  catch(e) { toast(e.message,'err'); }
}

// ── Init ──────────────────────────────────────────────────────────────────────
(async () => {
  try {
    const r = await api.get('/ping');
    if (r.authenticated) {
      g('login-page').classList.add('hidden');
      g('app').classList.remove('hidden');
      loadStatus(); loadProfiles();
      setInterval(loadStatus, 5000);
    }
  } catch(e) {}
})();

g('lp').addEventListener('keydown', e => { if(e.key==='Enter') doLogin(); });
</script>
</body>
</html>`
