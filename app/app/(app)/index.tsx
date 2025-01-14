import React, { useCallback, useEffect, useState } from "react";
import {
  StyleSheet,
  View,
  Text,
  Button,
  FlatList,
  StatusBar,
  Pressable,
  Alert,
  TouchableOpacity,
} from "react-native";
import { useSession } from "@/contexts/authContext";
import { useUser } from "@/contexts/userContext";
import Feather from "@expo/vector-icons/Feather";
import api from "@/scripts/axios";
import { Link } from "expo-router";

type ChatRoom = {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
  created_by: number;
  user: {
    name: string;
  };
};

export default function Home() {
  const { signOut } = useSession();
  const { user } = useUser();
  const [roomList, setRoomList] = useState<ChatRoom[]>([]);

  const fetchChats = useCallback(async () => {
    try {
      const response = await api.get<ChatRoom[]>("/chat");
      setRoomList(response.data);
    } catch (err) {
      console.error("Erro ao buscar chats:", err);
    }
  }, []);

  useEffect(() => {
    fetchChats();
  }, [fetchChats]);

  const deleteChat = async (chatId: number) => {
    try {
      await api.delete<ChatRoom>(`/chat/${chatId}`);
      setRoomList(roomList.filter((room) => room.id !== chatId));
    } catch (error) {
      console.log(error);
    }
  };

  const confirmDeletion = (id: number) =>
    Alert.alert(
      "Deseja realmente apagar essa sala?",
      "Está ação é irreversível",
      [
        {
          text: "Cancelar",
          style: "cancel",
        },
        { text: "Confirmar", onPress: () => deleteChat(id) },
      ]
    );

  const renderItem = ({ item }: { item: ChatRoom }) => (
<Link
  href={{ pathname: "/room/[id]",  params: { id: item.id, name: item.name } }}
  asChild
>
      <TouchableOpacity activeOpacity={0.5} style={styles.row}>
        <Text>{item.name}</Text>
        {user.id === item.created_by ? (
          <Pressable
            onPress={() => confirmDeletion(item.id)}
            style={({ pressed }) => [
              {
                backgroundColor: pressed
                  ? "#ff000020"
                  : styles.row.backgroundColor,
              },
              styles.deleteBtn,
            ]}
          >
            <Feather
              style={styles.deleteIcon}
              name="trash"
              size={22}
              color="black"
            />
          </Pressable>
        ) : (
          <></>
        )}
      </TouchableOpacity>
    </Link>
  );

  return (
    <View style={styles.container}>
      <FlatList
        data={roomList}
        keyExtractor={(item) => item.id.toString()}
        renderItem={renderItem}
      />
      <Button title="Logout" onPress={() => signOut()} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginTop: StatusBar.currentHeight,
    flex: 1,
  },
  row: {
    backgroundColor: "#c1c4c9",
    padding: 20,
    margin: 3,
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
  },
  deleteIcon: {
    color: "#FF0000",
  },
  deleteBtn: {
    padding: 5,
    borderRadius: 100000,
  },
});
