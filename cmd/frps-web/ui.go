package main

const uiHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>frps-web</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
:root{
  --primary:#8b5cf6;--primary-light:#a78bfa;--primary-dark:#7c3aed;
  --success:#10b981;--danger:#ef4444;--warning:#f59e0b;--info:#3b82f6;
  --bg:#f1f5f9;--sidebar:#ffffff;--card:#ffffff;
  --text:#1e293b;--text-muted:#64748b;--border:#e2e8f0;
  --shadow:0 1px 3px rgba(0,0,0,.08),0 1px 2px rgba(0,0,0,.04);
  --shadow-md:0 4px 6px -1px rgba(0,0,0,.07),0 2px 4px -1px rgba(0,0,0,.04);
  --radius:12px;--radius-sm:8px;
}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:var(--bg);color:var(--text);min-height:100vh}

/* ── Layout ── */
#login-page{display:flex;align-items:center;justify-content:center;min-height:100vh;background:linear-gradient(135deg,#8b5cf6 0%,#6366f1 50%,#3b82f6 100%)}
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
.field input:focus{border-color:var(--primary);box-shadow:0 0 0 3px rgba(139,92,246,.1)}
.login-hint{font-size:12px;color:var(--text-muted);text-align:center;margin-top:16px}

/* ── Sidebar ── */
.sidebar-logo{padding:20px 16px;display:flex;align-items:center;gap:10px;border-bottom:1px solid var(--border)}
.sidebar-logo .logo-box{width:36px;height:36px;background:linear-gradient(135deg,var(--primary),var(--primary-dark));border-radius:8px;display:flex;align-items:center;justify-content:center;color:#fff;font-weight:800;font-size:11px;flex-shrink:0}
.sidebar-logo span{font-weight:700;font-size:16px;color:var(--text)}
.sidebar-logo small{font-size:11px;color:var(--text-muted);display:block}
.sidebar-nav{flex:1;padding:12px 8px;overflow-y:auto}
.nav-section{font-size:10px;font-weight:600;color:var(--text-muted);letter-spacing:.08em;text-transform:uppercase;padding:8px 8px 4px}
.nav-item{display:flex;align-items:center;gap:10px;padding:9px 12px;border-radius:var(--radius-sm);cursor:pointer;color:var(--text-muted);font-size:14px;font-weight:500;transition:.15s;margin-bottom:2px}
.nav-item:hover{background:#f8fafc;color:var(--text)}
.nav-item.active{background:linear-gradient(135deg,rgba(139,92,246,.12),rgba(139,92,246,.06));color:var(--primary);font-weight:600}
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
.status-dot{width:7px;height:7px;border-radius:50%}
.running .status-dot{background:#16a34a;animation:pulse 2s infinite}
.stopped .status-dot{background:#dc2626}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:.4}}
.ctrl-btns{display:flex;gap:6px}

/* ── Buttons ── */
.btn{display:inline-flex;align-items:center;gap:6px;padding:8px 16px;border-radius:var(--radius-sm);font-size:13px;font-weight:500;cursor:pointer;border:none;transition:.15s;white-space:nowrap}
.btn:disabled{opacity:.5;cursor:not-allowed}
.btn-primary{background:linear-gradient(135deg,var(--primary),var(--primary-dark));color:#fff;box-shadow:0 2px 8px rgba(139,92,246,.3)}
.btn-primary:hover:not(:disabled){box-shadow:0 4px 12px rgba(139,92,246,.4);transform:translateY(-1px)}
.btn-success{background:linear-gradient(135deg,#10b981,#059669);color:#fff;box-shadow:0 2px 8px rgba(16,185,129,.3)}
.btn-danger{background:linear-gradient(135deg,#ef4444,#dc2626);color:#fff;box-shadow:0 2px 8px rgba(239,68,68,.3)}
.btn-warning{background:linear-gradient(135deg,#f59e0b,#d97706);color:#fff}
.btn-ghost{background:transparent;color:var(--text-muted);border:1.5px solid var(--border)}
.btn-ghost:hover{background:#f8fafc;color:var(--text)}
.btn-sm{padding:5px 10px;font-size:12px}
.btn-icon{padding:7px;border-radius:var(--radius-sm);background:transparent;color:var(--text-muted);border:1.5px solid var(--border);cursor:pointer;display:flex;align-items:center;transition:.15s}
.btn-icon:hover{background:#f8fafc;color:var(--text)}
.btn-login{width:100%;padding:12px;background:linear-gradient(135deg,var(--primary),var(--primary-dark));color:#fff;border:none;border-radius:var(--radius-sm);font-size:15px;font-weight:600;cursor:pointer;transition:.2s;margin-top:8px}
.btn-login:hover{transform:translateY(-1px);box-shadow:0 8px 20px rgba(139,92,246,.35)}

/* ── Cards ── */
.card{background:var(--card);border-radius:var(--radius);box-shadow:var(--shadow);border:1px solid var(--border)}
.card-header{padding:16px 20px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between}
.card-title{font-size:14px;font-weight:600;color:var(--text);display:flex;align-items:center;gap:8px}
.card-body{padding:20px}

/* ── Stat cards ── */
.stats-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:16px;margin-bottom:24px}
.stat-card{background:var(--card);border-radius:var(--radius);padding:20px;border:1px solid var(--border);position:relative;overflow:hidden}
.stat-card::before{content:'';position:absolute;top:0;left:0;right:0;height:3px}
.stat-card.purple::before{background:linear-gradient(90deg,#8b5cf6,#a78bfa)}
.stat-card.blue::before{background:linear-gradient(90deg,#3b82f6,#60a5fa)}
.stat-card.green::before{background:linear-gradient(90deg,#10b981,#34d399)}
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
.badge-online{background:#dcfce7;color:#166534}
.badge-blocked{background:#fef2f2;color:#991b1b}

/* ── Forms ── */
.form-grid{display:grid;grid-template-columns:1fr 1fr;gap:16px}
.form-field{display:flex;flex-direction:column;gap:5px}
.form-field.full{grid-column:1/-1}
.form-field label{font-size:12px;font-weight:600;color:var(--text-muted);text-transform:uppercase;letter-spacing:.05em}
.form-field input,.form-field select{padding:9px 12px;border:1.5px solid var(--border);border-radius:var(--radius-sm);font-size:13px;color:var(--text);outline:none;transition:.2s;background:#fff}
.form-field input:focus,.form-field select:focus{border-color:var(--primary);box-shadow:0 0 0 3px rgba(139,92,246,.1)}

/* ── Log viewer ── */
.log-box{background:#0f172a;border-radius:var(--radius-sm);padding:14px 16px;font-family:monospace;font-size:12px;line-height:1.6;height:420px;overflow-y:auto;color:#94a3b8}
.log-box .log-info{color:#94a3b8}
.log-box .log-warn{color:#fbbf24}
.log-box .log-error{color:#f87171}
.log-box .log-debug{color:#6b7280}

/* ── Traffic ── */
.traffic-bar{height:6px;background:#f1f5f9;border-radius:3px;overflow:hidden;margin-top:4px}
.traffic-bar-inner{height:100%;border-radius:3px;background:linear-gradient(90deg,var(--primary),var(--primary-light))}

/* ── Modal ── */
.modal-bg{position:fixed;inset:0;background:rgba(0,0,0,.4);backdrop-filter:blur(4px);z-index:1000;display:flex;align-items:center;justify-content:center;opacity:0;pointer-events:none;transition:.2s}
.modal-bg.open{opacity:1;pointer-events:all}
.modal{background:#fff;border-radius:16px;width:520px;max-width:95vw;max-height:90vh;overflow-y:auto;box-shadow:0 25px 50px rgba(0,0,0,.2);transform:scale(.95);transition:.2s}
.modal-bg.open .modal{transform:scale(1)}
.modal-head{padding:20px 24px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;position:sticky;top:0;background:#fff;z-index:1}
.modal-title{font-size:16px;font-weight:600}
.modal-body{padding:24px}
.modal-foot{padding:16px 24px;border-top:1px solid var(--border);display:flex;justify-content:flex-end;gap:10px;position:sticky;bottom:0;background:#fff}
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

/* ── Misc ── */
.hidden{display:none!important}
.flex{display:flex}.items-center{align-items:center}.gap-2{gap:8px}.gap-3{gap:12px}.justify-between{justify-content:space-between}
.mb-4{margin-bottom:16px}.mb-6{margin-bottom:24px}
.page{display:none}.page.active{display:block}
.section-title{font-size:14px;font-weight:600;color:var(--text);margin-bottom:16px;display:flex;align-items:center;gap:8px}
.font-mono{font-family:monospace}
.username-chip{display:inline-flex;align-items:center;gap:5px;padding:3px 10px;background:#ede9fe;color:#7c3aed;border-radius:20px;font-size:12px;font-weight:600}
.port-chip{display:inline-flex;padding:2px 7px;background:#dbeafe;color:#1d4ed8;border-radius:6px;font-size:11px;font-family:monospace;margin:1px}
</style>
</head>
<body>

<div id="toast"></div>

<!-- Login -->
<div id="login-page">
  <div class="login-card">
    <div class="login-logo">
      <div class="logo-box">FRP</div>
      <h1>frps-web</h1>
    </div>
    <h2>服务端管理面板</h2>
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
      <div><span>frps-web</span><small>服务端管理</small></div>
    </div>
    <nav class="sidebar-nav">
      <div class="nav-section">监控</div>
      <div class="nav-item active" data-page="overview" onclick="nav('overview')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
        概览
      </div>
      <div class="nav-item" data-page="clients" onclick="nav('clients')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><circle cx="9" cy="7" r="4"/><path d="M3 21v-2a4 4 0 014-4h4a4 4 0 014 4v2"/><path d="M16 3.13a4 4 0 010 7.75"/><path d="M21 21v-2a4 4 0 00-3-3.85"/></svg>
        连接设备
      </div>
      <div class="nav-item" data-page="traffic" onclick="nav('traffic')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
        代理流量
      </div>
      <div class="nav-item" data-page="logs" onclick="nav('logs')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
        连接日志
      </div>
      <div class="nav-section">配置</div>
      <div class="nav-item" data-page="settings" onclick="nav('settings')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><circle cx="12" cy="12" r="3"/><path d="M19.07 4.93l-1.41 1.41M4.93 4.93l1.41 1.41M12 2v2M12 20v2M20 12h2M2 12h2M19.07 19.07l-1.41-1.41M4.93 19.07l1.41-1.41"/></svg>
        服务设置
      </div>
      <div class="nav-section">账户</div>
      <div class="nav-item" data-page="account" onclick="nav('account')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/></svg>
        账户安全
      </div>
    </nav>
    <div class="sidebar-footer">
      <div class="nav-item" onclick="doLogout()">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
        退出登录
      </div>
    </div>
  </aside>

  <main class="main">
    <header class="topbar">
      <div class="topbar-left">
        <div>
          <div class="page-title" id="page-title">概览</div>
          <div style="font-size:12px;color:var(--text-muted)" id="server-sub">-</div>
        </div>
      </div>
      <div class="topbar-right">
        <div id="status-badge" class="status-badge stopped"><span class="status-dot"></span><span id="status-text">已停止</span></div>
        <div class="ctrl-btns">
          <button class="btn btn-success btn-sm" onclick="doStart()">
            <svg width="12" height="12" fill="currentColor" viewBox="0 0 24 24"><polygon points="5,3 19,12 5,21"/></svg> 启动
          </button>
          <button class="btn btn-danger btn-sm" onclick="doStop()">
            <svg width="12" height="12" fill="currentColor" viewBox="0 0 24 24"><rect x="3" y="3" width="18" height="18" rx="2"/></svg> 停止
          </button>
          <button class="btn btn-warning btn-sm" onclick="doRestart()">
            <svg width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg> 重启
          </button>
          <button class="btn-icon" onclick="refreshAll()" title="刷新">
            <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
          </button>
        </div>
        <span id="user-badge" style="font-size:12px;color:var(--text-muted)">admin</span>
      </div>
    </header>

    <div class="content">

      <!-- Overview -->
      <div id="page-overview" class="page active">
        <div class="stats-grid">
          <div class="stat-card purple"><div class="stat-label">在线客户端</div><div class="stat-value" id="s-clients">0</div></div>
          <div class="stat-card blue"><div class="stat-label">活跃代理</div><div class="stat-value" id="s-proxies">0</div></div>
          <div class="stat-card green"><div class="stat-label">今日流入</div><div class="stat-value" id="s-in">0 B</div></div>
          <div class="stat-card orange"><div class="stat-label">今日流出</div><div class="stat-value" id="s-out">0 B</div></div>
        </div>
        <div class="card">
          <div class="card-header">
            <div class="card-title">
              <svg width="15" height="15" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><circle cx="12" cy="12" r="3"/><path d="M12 1v4M12 19v4M4.22 4.22l2.83 2.83M16.95 16.95l2.83 2.83M1 12h4M19 12h4M4.22 19.78l2.83-2.83M16.95 7.05l2.83-2.83"/></svg>
              服务器运行信息
            </div>
          </div>
          <div class="card-body">
            <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:12px" id="srv-info-grid">
              <div><div style="font-size:11px;color:var(--text-muted);font-weight:600;text-transform:uppercase;letter-spacing:.05em">监听端口</div><div class="font-mono" style="margin-top:4px;font-weight:600" id="si-port">-</div></div>
              <div><div style="font-size:11px;color:var(--text-muted);font-weight:600;text-transform:uppercase;letter-spacing:.05em">运行时长</div><div class="font-mono" style="margin-top:4px;font-weight:600" id="si-uptime">-</div></div>
              <div><div style="font-size:11px;color:var(--text-muted);font-weight:600;text-transform:uppercase;letter-spacing:.05em">启动时间</div><div class="font-mono" style="margin-top:4px;font-weight:600" id="si-starttime">-</div></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Clients -->
      <div id="page-clients" class="page">
        <div class="card">
          <div class="card-header">
            <div class="card-title">
              <svg width="15" height="15" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><circle cx="9" cy="7" r="4"/><path d="M3 21v-2a4 4 0 014-4h4a4 4 0 014 4v2"/></svg>
              连接设备
            </div>
            <button class="btn-icon" onclick="loadClients()">
              <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
            </button>
          </div>
          <div class="table-wrap">
            <table>
              <thead><tr><th>用户名</th><th>主机名 / 系统</th><th>IP 地址</th><th>占用端口</th><th>连接时间</th><th>状态</th><th>操作</th></tr></thead>
              <tbody id="clients-body"><tr class="empty-row"><td colspan="7">暂无连接设备</td></tr></tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Traffic -->
      <div id="page-traffic" class="page">
        <div class="card">
          <div class="card-header">
            <div class="card-title">
              <svg width="15" height="15" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
              代理流量
            </div>
            <button class="btn-icon" onclick="loadTraffic()">
              <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
            </button>
          </div>
          <div class="table-wrap">
            <table>
              <thead><tr><th>代理名称</th><th>类型</th><th>客户端</th><th>今日流入</th><th>今日流出</th><th>当前连接</th><th>最后启动</th></tr></thead>
              <tbody id="traffic-body"><tr class="empty-row"><td colspan="7">暂无数据</td></tr></tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Logs -->
      <div id="page-logs" class="page">
        <div class="card">
          <div class="card-header">
            <div class="card-title">
              <svg width="15" height="15" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
              实时日志
            </div>
            <div class="flex gap-2">
              <label style="display:flex;align-items:center;gap:6px;font-size:12px;color:var(--text-muted);cursor:pointer">
                <input type="checkbox" id="log-auto" checked> 自动滚动
              </label>
              <button class="btn btn-ghost btn-sm" onclick="g('log-box').innerHTML='';logLines=[]">清空</button>
            </div>
          </div>
          <div class="card-body" style="padding:0">
            <div class="log-box" id="log-box"></div>
          </div>
        </div>
      </div>

      <!-- Settings -->
      <div id="page-settings" class="page">
        <div class="card" style="max-width:600px">
          <div class="card-header"><div class="card-title">frps 服务配置</div></div>
          <div class="card-body">
            <div class="form-grid" style="gap:16px">
              <div class="form-field"><label>监听端口 BindPort</label><input id="cfg-port" type="number" placeholder="7000"></div>
              <div class="form-field"><label>认证 Token</label><input id="cfg-token" type="password" placeholder="留空则不验证"></div>
              <div class="form-field"><label>HTTP 虚拟主机端口</label><input id="cfg-http" type="number" placeholder="0 = 禁用"></div>
              <div class="form-field"><label>HTTPS 虚拟主机端口</label><input id="cfg-https" type="number" placeholder="0 = 禁用"></div>
              <div class="form-field full"><label>子域名根域 SubdomainHost</label><input id="cfg-subdomain" type="text" placeholder="例如: frps.example.com"></div>
              <div class="form-field"><label>日志级别</label>
                <select id="cfg-loglevel">
                  <option value="info">info</option><option value="debug">debug</option>
                  <option value="warn">warn</option><option value="error">error</option>
                </select>
              </div>
            </div>
            <div style="margin-top:20px;display:flex;gap:10px">
              <button class="btn btn-primary" onclick="saveSettings()">保存设置</button>
              <button class="btn btn-ghost" onclick="loadSettings()">重置</button>
            </div>
            <p style="font-size:12px;color:var(--text-muted);margin-top:12px">保存后点击「重启」使配置生效</p>
          </div>
        </div>
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

<script>
const g = id => document.getElementById(id);
const esc = s => String(s ?? '').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');

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
};

function toast(msg, type='ok') {
  const el = document.createElement('div');
  el.className = 'toast-item toast-' + type;
  el.innerHTML = '<span>' + esc(msg) + '</span>';
  g('toast').appendChild(el);
  setTimeout(() => { el.style.animation='fadeOut .3s ease forwards'; setTimeout(() => el.remove(), 300); }, 3000);
}

function fmtBytes(n) {
  n = Number(n) || 0;
  if (n < 1024) return n + ' B';
  if (n < 1048576) return (n/1024).toFixed(1) + ' KB';
  if (n < 1073741824) return (n/1048576).toFixed(1) + ' MB';
  return (n/1073741824).toFixed(2) + ' GB';
}

// ── Auth ─────────────────────────────────────────────────────────────────────
async function doLogin() {
  const u = g('lu').value.trim(), p = g('lp').value;
  if (!u || !p) return;
  try {
    await api.post('/login', {Username: u, Password: p});
    g('login-page').classList.add('hidden');
    g('app').classList.remove('hidden');
    g('user-badge').textContent = u;
    refreshAll();
    setInterval(loadStatus, 5000);
    startLogPoll();
  } catch(e) { toast(e.message, 'err'); }
}

async function doLogout() {
  await api.post('/logout').catch(()=>{});
  g('app').classList.add('hidden');
  g('login-page').classList.remove('hidden');
  g('lp').value = '';
}

// ── Navigation ───────────────────────────────────────────────────────────────
const pageTitles = {overview:'概览', clients:'连接设备', traffic:'代理流量', logs:'连接日志', settings:'服务设置', account:'账户安全'};

function nav(page) {
  document.querySelectorAll('.page').forEach(el => el.classList.remove('active'));
  document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));
  g('page-' + page).classList.add('active');
  document.querySelector('.nav-item[data-page="' + page + '"]').classList.add('active');
  g('page-title').textContent = pageTitles[page] || page;
  if (page === 'clients') loadClients();
  if (page === 'traffic') loadTraffic();
  if (page === 'settings') loadSettings();
}

function refreshAll() { loadStatus(); }

// ── Status ───────────────────────────────────────────────────────────────────
async function loadStatus() {
  try {
    const d = await api.get('/frps/status');
    const badge = g('status-badge');
    const txt = g('status-text');
    if (d.running) {
      badge.className = 'status-badge running'; txt.textContent = '运行中';
    } else {
      badge.className = 'status-badge stopped'; txt.textContent = '已停止';
    }
    g('server-sub').textContent = d.running ? '端口 ' + (d.bindPort||7000) + ' · 运行中' : '未运行';
    g('s-clients').textContent = (d.stats && d.stats.totalClientConns) || 0;
    g('s-proxies').textContent = (d.stats && d.stats.totalProxyConns) || 0;
    g('s-in').textContent = fmtBytes(d.stats && d.stats.totalTrafficIn);
    g('s-out').textContent = fmtBytes(d.stats && d.stats.totalTrafficOut);
    g('si-port').textContent = d.bindPort || '-';
    g('si-uptime').textContent = d.uptime || '-';
    g('si-starttime').textContent = d.startTime ? new Date(d.startTime).toLocaleString() : '-';
  } catch(e) {}
}

// ── frps control ─────────────────────────────────────────────────────────────
async function doStart() {
  try { await api.post('/frps/start'); toast('启动中...','info'); setTimeout(loadStatus,1000); } catch(e) { toast(e.message,'err'); }
}
async function doStop() {
  try { await api.post('/frps/stop'); toast('已停止'); setTimeout(loadStatus,500); } catch(e) { toast(e.message,'err'); }
}
async function doRestart() {
  try { await api.post('/frps/restart'); toast('重启中...','info'); setTimeout(loadStatus,1500); } catch(e) { toast(e.message,'err'); }
}

// ── Clients ──────────────────────────────────────────────────────────────────
async function loadClients() {
  try {
    const list = await api.get('/clients');
    const tb = g('clients-body');
    if (!list || !list.length) {
      tb.innerHTML = '<tr class="empty-row"><td colspan="7">暂无连接设备</td></tr>'; return;
    }
    tb.innerHTML = list.map(c => {
      const ports = (c.ports||[]).map(p => '<span class="port-chip">:'+p+'</span>').join('');
      const since = c.connectedAt ? new Date(c.connectedAt).toLocaleString() : '-';
      const blocked = c.blocked;
      const username = c.user || c.key || '-';
      return '<tr>' +
        '<td><span class="username-chip">' + esc(username) + '</span></td>' +
        '<td>' + esc(c.hostname||'-') + '</td>' +
        '<td class="font-mono">' + esc(c.remoteIP||'-') + '</td>' +
        '<td>' + (ports||'-') + '</td>' +
        '<td class="font-mono">' + since + '</td>' +
        '<td><span class="badge ' + (blocked?'badge-blocked':'badge-online') + '">' + (blocked?'已拦截':'在线') + '</span></td>' +
        '<td><button class="btn btn-sm ' + (blocked?'btn-success':'btn-danger') + '" onclick="toggleBlock(\'' + esc(c.key) + '\',' + blocked + ')">' + (blocked?'解除拦截':'拦截') + '</button></td>' +
        '</tr>';
    }).join('');
  } catch(e) { toast(e.message,'err'); }
}

async function toggleBlock(key, blocked) {
  try {
    if (blocked) { await api.post('/clients/' + encodeURIComponent(key) + '/unblock'); toast('已解除拦截'); }
    else { await api.post('/clients/' + encodeURIComponent(key) + '/block'); toast('已拦截'); }
    loadClients();
  } catch(e) { toast(e.message,'err'); }
}

// ── Traffic ──────────────────────────────────────────────────────────────────
async function loadTraffic() {
  try {
    const list = await api.get('/proxies');
    const tb = g('traffic-body');
    const allProxies = [];
    if (list) {
      const types = ['tcp','udp','http','https','stcp','xtcp','tcpmux','sudp'];
      types.forEach(t => { if (list[t]) list[t].forEach(p => allProxies.push({...p, ptype:t})); });
    }
    if (!allProxies.length) {
      tb.innerHTML = '<tr class="empty-row"><td colspan="7">暂无数据</td></tr>'; return;
    }
    tb.innerHTML = allProxies.map(p => {
      const tin = fmtBytes(p.todayTrafficIn);
      const tout = fmtBytes(p.todayTrafficOut);
      const last = p.lastStartTime ? new Date(p.lastStartTime*1000).toLocaleString() : '-';
      return '<tr>' +
        '<td class="font-mono">' + esc(p.name) + '</td>' +
        '<td><span class="badge badge-' + p.ptype + '">' + p.ptype.toUpperCase() + '</span></td>' +
        '<td>' + (p.clientVersion ? '<span class="username-chip">' + esc(p.clientVersion||p.runID||'-') + '</span>' : '-') + '</td>' +
        '<td>' + tin + '</td>' +
        '<td>' + tout + '</td>' +
        '<td>' + (p.curConns||0) + '</td>' +
        '<td class="font-mono" style="font-size:11px">' + last + '</td>' +
        '</tr>';
    }).join('');
  } catch(e) { toast(e.message,'err'); }
}

// ── Logs ─────────────────────────────────────────────────────────────────────
let logLines = [];
let lastLogCount = 0;
let logTimer = null;

function startLogPoll() {
  if (logTimer) clearInterval(logTimer);
  logTimer = setInterval(loadLogs, 2500);
}

async function loadLogs() {
  try {
    const logs = await api.get('/logs');
    if (!logs || logs.length === lastLogCount) return;
    lastLogCount = logs.length;
    const box = g('log-box');
    box.innerHTML = logs.map(l => {
      const cls = l.level === 'warn' ? 'log-warn' : l.level === 'error' ? 'log-error' : l.level === 'debug' ? 'log-debug' : 'log-info';
      return '<div class="' + cls + '">' + esc(l.time) + ' [' + (l.level||'I').toUpperCase()[0] + '] ' + esc(l.msg) + '</div>';
    }).join('');
    if (g('log-auto').checked) box.scrollTop = box.scrollHeight;
  } catch(e) {}
}

// ── Settings ─────────────────────────────────────────────────────────────────
async function loadSettings() {
  try {
    const d = await api.get('/settings');
    g('cfg-port').value = d.BindPort || d.bindPort || 7000;
    g('cfg-token').value = '';
    g('cfg-http').value = d.VhostHTTPPort || d.vhostHTTPPort || 0;
    g('cfg-https').value = d.VhostHTTPSPort || d.vhostHTTPSPort || 0;
    g('cfg-subdomain').value = d.SubDomainHost || d.subDomainHost || '';
    g('cfg-loglevel').value = d.LogLevel || d.logLevel || 'info';
  } catch(e) { toast(e.message,'err'); }
}

async function saveSettings() {
  try {
    const body = {
      BindPort: parseInt(g('cfg-port').value) || 7000,
      AuthToken: g('cfg-token').value,
      VhostHTTPPort: parseInt(g('cfg-http').value) || 0,
      VhostHTTPSPort: parseInt(g('cfg-https').value) || 0,
      SubDomainHost: g('cfg-subdomain').value.trim(),
      LogLevel: g('cfg-loglevel').value,
    };
    await api.put('/settings', body);
    toast('已保存，重启后生效');
  } catch(e) { toast(e.message,'err'); }
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
      refreshAll();
      setInterval(loadStatus, 5000);
      startLogPoll();
    }
  } catch(e) {}
})();

g('lp').addEventListener('keydown', e => { if(e.key==='Enter') doLogin(); });
</script>
</body>
</html>`
