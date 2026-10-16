import React, { useEffect, useState, useRef } from 'react';
import { aiApi } from '@/api/ai';
import { AIConversation, AIMessage } from '@/types/api';
import { formatJakartaDate, formatJakartaTimeOnly } from '@/utils/date';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { Skeleton } from '@/components/ui/EmptyState';
import {
  Sparkles,
  Send,
  Plus,
  Trash2,
  Bot,
  User,
  Wrench,
  MessageSquare,
  AlertCircle,
  HelpCircle,
} from 'lucide-react';

const SUGGESTED_PROMPTS = [
  'Berapa omzet laundry saya bulan ini?',
  'Berapa banyak pesanan yang belum selesai?',
  'Siapa saja pelanggan yang belum mengambil cucian?',
  'Layanan apa yang paling banyak dipesan minggu ini?',
  'Berapa total pengeluaran untuk deterjen bulan ini?',
  'Analisis kondisi bisnis laundry saya dan berikan rekomendasi.',
];

export const AIAssistantPage: React.FC = () => {
  const [conversations, setConversations] = useState<AIConversation[]>([]);
  const [activeConversationId, setActiveConversationId] = useState<number | null>(null);
  const [messages, setMessages] = useState<AIMessage[]>([]);

  const [inputMessage, setInputMessage] = useState('');
  const [isSending, setIsSending] = useState(false);
  const [isLoadingHistory, setIsLoadingHistory] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, isSending]);

  // Load conversations list
  const loadConversations = async () => {
    try {
      const res = await aiApi.getConversations(1, 30);
      setConversations(res.conversations);
      if (res.conversations.length > 0 && !activeConversationId) {
        // Load first conversation detail
        loadConversationDetail(res.conversations[0].id);
      }
    } catch (err: any) {
      setError(err.message || 'Gagal memuat percakapan AI.');
    } finally {
      setIsLoadingHistory(false);
    }
  };

  useEffect(() => {
    loadConversations();
  }, []);

  // Load active conversation messages
  const loadConversationDetail = async (id: number) => {
    setActiveConversationId(id);
    try {
      const detail = await aiApi.getConversationDetail(id);
      setMessages(detail.messages || []);
    } catch {
      setMessages([]);
    }
  };

  // Start fresh conversation
  const handleNewChat = () => {
    setActiveConversationId(null);
    setMessages([]);
    setInputMessage('');
  };

  // Delete conversation
  const handleDeleteConversation = async (id: number, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!window.confirm('Hapus riwayat percakapan ini?')) return;

    try {
      await aiApi.deleteConversation(id);
      setConversations(prev => prev.filter(c => c.id !== id));
      if (activeConversationId === id) {
        handleNewChat();
      }
    } catch (err: any) {
      alert(err.message || 'Gagal menghapus percakapan.');
    }
  };

  // Send prompt
  const handleSendMessage = async (textToSend?: string) => {
    const text = (textToSend || inputMessage).trim();
    if (!text || isSending) return;

    setInputMessage('');
    setError(null);

    // Optimistically append user message
    const tempUserMsg: AIMessage = {
      id: Date.now(),
      role: 'user',
      content: text,
      created_at: new Date().toISOString(),
    };
    setMessages(prev => [...prev, tempUserMsg]);
    setIsSending(true);

    try {
      const res = await aiApi.chat(text, activeConversationId || undefined);

      // Set active conversation ID if it was newly created
      if (!activeConversationId && res.conversation_id) {
        setActiveConversationId(res.conversation_id);
      }

      // Append model response
      const modelMsg: AIMessage = {
        id: Date.now() + 1,
        role: 'model',
        content: res.answer || res.message || '',
        tool_calls_json: res.tool_calls ? JSON.stringify(res.tool_calls) : undefined,
        created_at: new Date().toISOString(),
      };
      setMessages(prev => [...prev, modelMsg]);

      // Reload conversations list to update titles/ordering
      const convList = await aiApi.getConversations(1, 30);
      setConversations(convList.conversations);
    } catch (err: any) {
      setError(err.message || 'AI Assistant tidak dapat merespons saat ini.');
    } finally {
      setIsSending(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <div className="h-[calc(100vh-120px)] flex flex-col space-y-4">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 shrink-0">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-bold text-slate-900 tracking-tight">AI Business Assistant</h1>
            <span className="text-[10px] font-bold bg-teal-50 text-teal-800 border border-teal-200 px-2 py-0.5 rounded-full uppercase">
              Google Gemini
            </span>
          </div>
          <p className="text-xs text-slate-500 mt-0.5">
            Tanyakan data bisnis, tren pesanan, pelanggan, dan pengeluaran menggunakan bahasa Indonesia.
          </p>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={handleNewChat}
          leftIcon={<Plus className="w-4 h-4" />}
        >
          Percakapan Baru
        </Button>
      </div>

      {/* Main Chat Layout (Left: Conversations List, Right: Chat Interface) */}
      <div className="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-4 min-h-0">
        {/* Left: Conversations Sidebar (3 Cols) */}
        <div className="hidden lg:flex lg:col-span-4 flex-col bg-white rounded-xl border border-slate-200/80 overflow-hidden shadow-xs">
          <div className="p-3 border-b border-slate-100 flex items-center justify-between">
            <span className="text-xs font-bold text-slate-700">Riwayat Percakapan</span>
            <Button size="sm" variant="ghost" onClick={handleNewChat}>
              <Plus className="w-4 h-4" />
            </Button>
          </div>

          <div className="flex-1 overflow-y-auto p-2 space-y-1">
            {isLoadingHistory ? (
              Array.from({ length: 4 }).map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))
            ) : conversations.length > 0 ? (
              conversations.map(conv => (
                <div
                  key={conv.id}
                  onClick={() => loadConversationDetail(conv.id)}
                  className={`group p-2.5 rounded-lg text-xs cursor-pointer flex items-center justify-between transition-colors ${
                    activeConversationId === conv.id
                      ? 'bg-teal-50 text-teal-900 font-semibold border border-teal-200'
                      : 'text-slate-700 hover:bg-slate-50 border border-transparent'
                  }`}
                >
                  <div className="flex items-center gap-2 min-w-0 pr-2">
                    <MessageSquare className="w-4 h-4 text-slate-400 shrink-0" />
                    <span className="truncate">{conv.title || 'Percakapan AI'}</span>
                  </div>
                  <button
                    onClick={e => handleDeleteConversation(conv.id, e)}
                    className="opacity-0 group-hover:opacity-100 p-1 text-slate-400 hover:text-rose-600 rounded transition-opacity"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              ))
            ) : (
              <p className="text-xs text-slate-400 text-center py-8">
                Belum ada riwayat percakapan.
              </p>
            )}
          </div>
        </div>

        {/* Right: Active Chat Thread (8 Cols) */}
        <div className="lg:col-span-8 flex flex-col bg-white rounded-xl border border-slate-200/80 overflow-hidden shadow-xs">
          {/* Thread Message Container */}
          <div className="flex-1 overflow-y-auto p-4 sm:p-6 space-y-4">
            {messages.length === 0 ? (
              <div className="h-full flex flex-col items-center justify-center text-center p-6 space-y-4">
                <div className="w-12 h-12 rounded-2xl bg-teal-50 text-teal-700 flex items-center justify-center shadow-xs">
                  <Sparkles className="w-6 h-6" />
                </div>
                <div>
                  <h3 className="text-base font-bold text-slate-900">
                    Halo, Ada yang bisa saya bantu terkait bisnis laundry Anda?
                  </h3>
                  <p className="text-xs text-slate-500 max-w-md mt-1">
                    AI membaca data aktual langsung dari sistem SatSet Anda secara read-only.
                  </p>
                </div>

                {/* Prompt Pills */}
                <div className="pt-2 max-w-lg w-full space-y-2">
                  <p className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">
                    Contoh Pertanyaan Populer:
                  </p>
                  <div className="flex flex-wrap justify-center gap-2">
                    {SUGGESTED_PROMPTS.map((prompt, idx) => (
                      <button
                        key={idx}
                        onClick={() => handleSendMessage(prompt)}
                        className="text-xs text-slate-700 bg-slate-50 hover:bg-teal-50 hover:text-teal-800 hover:border-teal-300 border border-slate-200 px-3 py-1.5 rounded-full transition-all text-left"
                      >
                        {prompt}
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            ) : (
              messages.map(msg => {
                const isUser = msg.role === 'user';
                let parsedTools: any[] = [];
                if (msg.tool_calls_json) {
                  try {
                    parsedTools = JSON.parse(msg.tool_calls_json);
                  } catch {}
                }

                return (
                  <div
                    key={msg.id}
                    className={`flex items-start gap-3 ${
                      isUser ? 'flex-row-reverse' : 'flex-row'
                    }`}
                  >
                    <div
                      className={`w-8 h-8 rounded-full flex items-center justify-center shrink-0 text-xs font-semibold ${
                        isUser
                          ? 'bg-slate-900 text-white'
                          : 'bg-teal-700 text-white'
                      }`}
                    >
                      {isUser ? <User className="w-4 h-4" /> : <Bot className="w-4 h-4" />}
                    </div>

                    <div
                      className={`max-w-[85%] rounded-2xl px-4 py-3 text-xs leading-relaxed ${
                        isUser
                          ? 'bg-teal-700 text-white rounded-tr-none'
                          : 'bg-slate-50 border border-slate-200 text-slate-800 rounded-tl-none'
                      }`}
                    >
                      {/* Tool Calls Badge */}
                      {parsedTools.length > 0 && (
                        <div className="mb-2 pb-2 border-b border-slate-200 flex flex-wrap gap-1.5">
                          {parsedTools.map((tc, idx) => (
                            <span
                              key={idx}
                              className="inline-flex items-center gap-1 text-[10px] bg-white border border-slate-200 text-slate-600 px-2 py-0.5 rounded-full font-medium"
                            >
                              <Wrench className="w-2.5 h-2.5 text-teal-600" />
                              Tool: {tc.name}
                            </span>
                          ))}
                        </div>
                      )}

                      {/* Content text */}
                      <div className="whitespace-pre-wrap">{msg.content}</div>

                      <span
                        className={`block text-[10px] mt-1.5 ${
                          isUser ? 'text-teal-200 text-right' : 'text-slate-400'
                        }`}
                      >
                        {formatJakartaTimeOnly(msg.created_at)}
                      </span>
                    </div>
                  </div>
                );
              })
            )}

            {/* Loading Indicator */}
            {isSending && (
              <div className="flex items-start gap-3">
                <div className="w-8 h-8 rounded-full bg-teal-700 text-white flex items-center justify-center shrink-0 text-xs font-semibold">
                  <Bot className="w-4 h-4" />
                </div>
                <div className="bg-slate-50 border border-slate-200 rounded-2xl rounded-tl-none px-4 py-3 text-xs text-slate-500 flex items-center gap-2">
                  <Sparkles className="w-4 h-4 text-teal-600 animate-spin" />
                  <span>Gemini sedang menganalisis data laundry Anda...</span>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* Error Banner */}
          {error && (
            <div className="px-4 py-2 bg-rose-50 border-t border-rose-200 text-xs text-rose-800 flex items-center gap-2">
              <AlertCircle className="w-4 h-4 text-rose-600 shrink-0" />
              <span>{error}</span>
            </div>
          )}

          {/* Input Prompt Box */}
          <div className="p-3 border-t border-slate-200/80 bg-white">
            <div className="flex items-end gap-2 bg-slate-50 border border-slate-200 rounded-xl p-2 focus-within:ring-2 focus-within:ring-teal-600 focus-within:bg-white transition-all">
              <textarea
                value={inputMessage}
                onChange={e => setInputMessage(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Ketik pertanyaan untuk AI Assistant... (Tekan Enter untuk mengirim)"
                rows={2}
                className="flex-1 bg-transparent text-xs text-slate-900 placeholder:text-slate-400 focus:outline-none resize-none px-1 py-0.5"
              />
              <Button
                variant="primary"
                size="sm"
                onClick={() => handleSendMessage()}
                disabled={!inputMessage.trim() || isSending}
                isLoading={isSending}
                className="shrink-0 h-9 px-3 rounded-lg"
              >
                <Send className="w-4 h-4" />
              </Button>
            </div>
            <div className="mt-1.5 flex items-center justify-between text-[11px] text-slate-400 px-1">
              <span>Data bersifat read-only • Privasi terjaga di server internal</span>
              <span className="hidden sm:inline">Shift + Enter untuk baris baru</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
