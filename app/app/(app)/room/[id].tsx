import {
  ActivityIndicator,
  FlatList,
  StyleSheet,
  TextInput,
  TouchableOpacity,
  View,
} from "react-native";
import { useLocalSearchParams, useNavigation } from "expo-router";
import { useEffect, useRef, useState } from "react";
import api from "@/scripts/axios";
import { useUser } from "@/contexts/userContext";
import Feather from "@expo/vector-icons/Feather";
import MessageItem from "@/components/message";
import { useSession } from "@/contexts/authContext";

interface Message {
  id: number;
  chat_id: number;
  content: string;
  created_at: string;
  user: {
    name: string;
  };
  user_id: number | null; // Permitir null em user_id
}

export default function Room() {
  const { user } = useUser();
  const { session } = useSession();
  const local = useLocalSearchParams();
  const navigation = useNavigation();
  const roomId = local.id;

  const [messages, setMessages] = useState<Message[]>([]);
  const [messageInput, setMessageInput] = useState("");
  const [page, setPage] = useState(0);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);

  const getMessages = async (pageNumber: number) => {
    if (loading || !hasMore) return;

    setLoading(true);

    try {
      const response = await api.get(`/messages/${roomId}/${pageNumber}`);
      if (response.data.length === 0) {
        setHasMore(false);
        return;
      }

      setMessages((prevMessages) => [...prevMessages, ...response.data]);
      setPage(pageNumber + 1);
    } catch (error) {
      console.log(error);
    } finally {
      setLoading(false);
    }
  };

  const connectWebSocket = () => {
    const baseUrl = process.env.EXPO_PUBLIC_API_URL;
    const socketUrl = `${baseUrl}chat/${roomId}?token=${session}`;
    socketRef.current = new WebSocket(socketUrl);

    socketRef.current.onopen = () => {
      setIsConnected(true);
    };

    socketRef.current.onmessage = (event) => {
      const newMessage = JSON.parse(event.data);
      console.log("Mensagem recebida do WebSocket:", newMessage);
      const obj = {
        id: Date.now(),
        content: newMessage.content,
        created_at: newMessage.created_at,
        user_id: newMessage.user_id,
        chat_id: newMessage.chat_id,
        user: {
          name: newMessage.user.name,
        },
      };
      setMessages((prevMessages) => [obj, ...prevMessages]);
    };

    socketRef.current.onerror = (error) => {
      console.error("Erro na conexão WebSocket:", error);
    };

    socketRef.current.onclose = () => {
      setIsConnected(false);
    };
  };

  const sendMessage = () => {
    if (socketRef.current && messageInput.trim()) {
      const obj: Message = {
        // Declarar explicitamente o tipo Message
        id: Date.now(),
        content: messageInput,
        created_at: new Date().toISOString(),
        user_id: user.id,
        chat_id: Number(roomId),
        user: {
          name: "Eu",
        },
      };

      socketRef.current.send(JSON.stringify(obj));
      setMessages((prevMessages) => [obj, ...prevMessages]); // Sempre retorna um array de Message[]
      setMessageInput("");
    }
  };

  useEffect(() => {
    navigation.setOptions({ headerShown: true, title: local.name });
    getMessages(0);
    connectWebSocket();
  }, []);

  const renderItem = ({ item }: { item: Message }) => (
    <MessageItem item={item} isOwnMessage={user.id === item.user_id} />
  );

  return (
    <View style={styles.container}>
      <FlatList
        contentContainerStyle={styles.messageList}
        data={messages}
        keyExtractor={(item) => item.created_at}
        renderItem={renderItem}
        onEndReached={() => getMessages(page)}
        onEndReachedThreshold={0.5}
        ListFooterComponent={
          loading ? <ActivityIndicator size="large" color="#0000ff" /> : null
        }
        inverted={true}
        keyboardShouldPersistTaps="handled"
      />
      <View style={styles.messageInputWrapper}>
        <TextInput
          style={styles.messageInput}
          value={messageInput}
          onChangeText={setMessageInput}
          placeholder="Digite uma mensagem"
          multiline={true}
        />
        <TouchableOpacity onPress={() => sendMessage()} style={styles.sendIcon}>
          <Feather name="send" size={20} color="black" />
        </TouchableOpacity>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 10,
  },
  ownMessage: {
    alignSelf: "flex-end",
    backgroundColor: "#9BFFAE",
    padding: 10,
    borderRadius: 10,
    marginBottom: 5,
    maxWidth: "75%",
  },
  someoneMessage: {
    alignSelf: "flex-start",
    backgroundColor: "#B2B4B2",
    padding: 10,
    borderRadius: 10,
    marginBottom: 5,
    maxWidth: "75%",
  },
  messageInput: {
    width: "90%",
    paddingVertical: 10,
    borderRadius: 20,
  },
  messageInputWrapper: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    paddingHorizontal: 10,
    borderWidth: 0.5,
    borderRadius: 20,
    marginTop: 10,
    backgroundColor: "#fff",
    shadowColor: "#000",
    shadowOffset: {
      width: 0,
      height: 1,
    },
    shadowOpacity: 0.2,
    shadowRadius: 1.41,

    elevation: 2,
  },
  sendIcon: {
    position: "absolute",
    right: 10,
    top: "50%",
    transform: [{ translateY: -10 }],
  },
  messageList: {
    justifyContent: "flex-end",
    paddingBottom: 10,
  },
});
