<script>
  export let instances = [];
  export let minimalMode = false;
  export let onStart;
  export let onStop;
  export let onEdit;
  export let onDuplicate;
  export let onDelete;
  export let onCheckProxy;

  $: columns = minimalMode ? [
    { label: '名称', field: 'name', width: '13%' },
    { label: '状态', field: 'status', width: '10%' },
    { label: '标签', field: 'tags', width: '40%' },
    { label: '代理位置', field: 'proxy', width: '17%' },
    { label: '操作', field: 'actions', width: '20%' }
  ] : [
    { label: '名称', field: 'name', width: '20%' },
    { label: '状态', field: 'status', width: '10%' },
    { label: '标签', field: 'tags', width: '15%' },
    { label: '路径', field: 'path', width: '25%' },
    { label: '参数', field: 'args', width: '15%' },
    { label: '操作', field: 'actions', width: '15%' }
  ];

  function getLatencyColor(latency) {
    if (!latency) return 'var(--text-muted)';
    if (latency < 100) return 'var(--success-color)';
    if (latency < 300) return '#f59e0b'; // warning
    return 'var(--danger-color)';
  }
</script>

<div class="table-container">
  <table>
    <thead>
      <tr>
        {#each columns as col}
          <th style="width: {col.width}">{col.label}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each instances as instance (instance.id)}
        <tr>
          <td>
            <div class="name-cell">
              <strong>{instance.name}</strong>
              <div class="id-text" title={instance.id}>#{instance.id.slice(0, 4)}</div>
            </div>
          </td>
          <td>
            {#if instance.running}
              <span class="badge success">运行中</span>
            {:else}
              <span class="badge muted">已停止</span>
            {/if}
          </td>
          <td>
            {#if instance.tags && instance.tags.length > 0}
              <div class="tags-list">
                {#each instance.tags as tag}
                  <span class="tag-badge">{tag}</span>
                {/each}
              </div>
            {:else}
              <span class="text-muted">-</span>
            {/if}
          </td>
          {#if minimalMode}
            <td>
               <div style="display: flex; align-items: center; gap: 0.5rem; justify-content: space-between;">
                 <span 
                   title={instance.proxyDetail || '暂无详细信息'}
                   style="cursor: help; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 150px; display: inline-block;"
                 >
                   {#if instance.proxyRegion}
                     {instance.proxyRegion} <span style="color: {getLatencyColor(instance.proxyLatency)}; font-weight: 500;">{instance.proxyLatency}ms</span>
                   {:else}
                     -
                   {/if}
                 </span>
                 <div style="display: flex; gap: 0.25rem;">
                   <button class="btn-xs" title="检测 CN" on:click={() => onCheckProxy(instance.id, 'cn')}>CN</button>
                   <button class="btn-xs" title="检测 Global" on:click={() => onCheckProxy(instance.id, 'global')}>GL</button>
                 </div>
               </div>
            </td>
          {:else}
            <td>
              <div class="path-text" title={instance.path}>{instance.path}</div>
              <div class="path-subtext" title={instance.userDataDir}>{instance.userDataDir || '默认数据目录'}</div>
            </td>
            <td>
              <div class="args-text" title={instance.args ? instance.args.join(' ') : ''}>{instance.args ? instance.args.join(' ') : '-'}</div>
            </td>
          {/if}
          <td>
            <div class="actions">
              {#if !instance.running}
                <button class="btn-sm btn-primary" on:click={() => onStart(instance.id)}>启动</button>
              {:else}
                <button class="btn-sm btn-danger" on:click={() => onStop(instance.id)}>停止</button>
              {/if}
              <button class="icon-btn" title="编辑" on:click={() => onEdit(instance)}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
              </button>
              <button class="icon-btn" title="复制" on:click={() => onDuplicate(instance)}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
              </button>
              <button class="icon-btn danger" title="删除" on:click={() => onDelete(instance.id)}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
              </button>
            </div>
          </td>
        </tr>
      {/each}
      {#if instances.length === 0}
        <tr>
          <td colspan="6" style="text-align: center; padding: 2rem; color: var(--text-muted);">
            暂无数据
          </td>
        </tr>
      {/if}
    </tbody>
  </table>
</div>

<style>
  .table-container {
    background: white;
    border-radius: 0.5rem;
    box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1);
    border: 1px solid var(--border-color);
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.875rem;
    text-align: left;
  }

  th {
    padding: 0.75rem 1rem;
    background-color: #f8fafc;
    color: var(--text-muted);
    font-weight: 600;
    border-bottom: 1px solid var(--border-color);
    white-space: nowrap;
  }

  td {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--border-color);
    vertical-align: middle;
  }

  tr:last-child td {
    border-bottom: none;
  }

  .name-cell {
    display: flex;
    flex-direction: column;
  }

  .id-text {
    font-size: 0.75rem;
    color: var(--text-muted);
  }

  .badge {
    display: inline-flex;
    padding: 0.125rem 0.5rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 500;
  }

  .badge.success {
    background-color: #dcfce7;
    color: #166534;
  }

  .badge.muted {
    background-color: #f1f5f9;
    color: #64748b;
  }

  .tags-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
  }

  .tag-badge {
    background-color: #e2e8f0;
    color: #475569;
    padding: 0.125rem 0.375rem;
    border-radius: 0.25rem;
    font-size: 0.75rem;
  }

  .text-muted {
    color: var(--text-muted);
  }

  .path-text {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 200px;
    font-family: monospace;
  }
  
  .path-subtext {
    font-size: 0.75rem;
    color: var(--text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 200px;
  }

  .args-text {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 150px;
    color: var(--text-muted);
    font-family: monospace;
    font-size: 0.75rem;
  }

  .actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .btn-sm {
    padding: 0.25rem 0.75rem;
    font-size: 0.75rem;
    border-radius: 0.25rem;
    border: none;
    cursor: pointer;
    font-weight: 500;
  }

  .btn-primary {
    background-color: var(--primary-color);
    color: white;
  }
  
  .btn-primary:active {
     opacity: 0.9;
  }

  .btn-danger {
    background-color: #fee2e2;
    color: #ef4444;
  }
  
  .btn-danger:active {
    background-color: #fecaca;
  }

  .icon-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    color: #64748b;
    padding: 0.375rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 0.25rem;
    transition: all 0.2s;
  }
  
  .icon-btn:hover {
    background-color: #f1f5f9;
    color: #334155;
  }

  .icon-btn.danger:hover {
    background-color: #fee2e2;
    color: #ef4444;
  }

  .btn-xs {
    padding: 0.125rem 0.375rem;
    font-size: 0.7rem;
    background-color: white;
    border: 1px solid var(--border-color);
    border-radius: 0.25rem;
    cursor: pointer;
    color: var(--primary-color);
    font-weight: 500;
  }
  .btn-xs:hover {
      background-color: #f8fafc;
  }
</style>
