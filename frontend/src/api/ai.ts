import { apiFetch } from './client';
import { AIChatResponse, AIConversation, PaginationMeta } from '@/types/api';

export const aiApi = {
  chat: async (message: string, conversationID?: number): Promise<AIChatResponse> => {
    const res = await apiFetch<AIChatResponse>('/ai/chat', {
      method: 'POST',
      body: JSON.stringify({
        question: message,
        message,
        conversation_id: conversationID || 0,
      }),
    });
    return res.data!;
  },

  getConversations: async (
    page: number = 1,
    perPage: number = 20
  ): Promise<{ conversations: AIConversation[]; meta?: PaginationMeta }> => {
    const res = await apiFetch<AIConversation[]>(
      `/ai/conversations?page=${page}&per_page=${perPage}`
    );
    return {
      conversations: res.data || [],
      meta: res.meta,
    };
  },

  getConversationDetail: async (id: number): Promise<AIConversation> => {
    const res = await apiFetch<AIConversation>(`/ai/conversations/${id}`);
    return res.data!;
  },

  deleteConversation: async (id: number): Promise<void> => {
    await apiFetch(`/ai/conversations/${id}`, {
      method: 'DELETE',
    });
  },

  listTools: async (): Promise<any[]> => {
    const res = await apiFetch<any[]>('/ai/tools');
    return res.data || [];
  },
};
