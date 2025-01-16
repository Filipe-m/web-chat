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
  RefreshControl,
} from "react-native";
import { useSession } from "@/contexts/authContext";
import { useUser } from "@/contexts/userContext";
import Feather from "@expo/vector-icons/Feather";
import api from "@/scripts/axios";
import { Link } from "expo-router";
import CreateChat from "@/components/createChat";
import DefaultBtn from "@/components/defaultButton";

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
  const [isModalVisible, setIsModalVisible] = useState<boolean>(false);
  const [refreshing, setRefreshing] = React.useState(false);

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

  const onModalClose = () => {
    setIsModalVisible(false);
  };

  const onModalOpen = () => {
    setIsModalVisible(true);
  };

  const renderItem = ({ item }: { item: ChatRoom }) => (
    <Link
      href={{
        pathname: "/room/[id]",
        params: { id: item.id, name: item.name },
      }}
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
          <Text></Text>
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
        contentContainerStyle={styles.flatList}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={fetchChats} />
        }
      />

      <View style={styles.buttonsContainer}>
        <DefaultBtn
          text="Criar sala"
          method={onModalOpen}
          containerStyle={{
            backgroundColor: "#FFFFFF",
            borderWidth: 1,
            borderColor: "#000000",
            width: "90%",
          }}
          pressedContainerStyle={{
            backgroundColor: "#808891",
            borderWidth: 1,
            width: "90%",
          }}
          textStyle={{
            fontSize: 16,
            color: "#000000",
          }}
          pressedTextStyle={{
            color: "#FFFFFF",
          }}
        />
        <DefaultBtn
          text="Logout"
          method={signOut}
          containerStyle={{
            backgroundColor: "#FFFFFF",
            borderWidth: 1,
            borderColor: "#000000",
            width: "90%",
          }}
          pressedContainerStyle={{
            backgroundColor: "#808891",
            borderWidth: 1,
            width: "90%",
          }}
          textStyle={{
            color: "#000000",
          }}
          pressedTextStyle={{
            color: "#FFFFFF",
          }}
        />
      </View>
      <CreateChat
        isVisible={isModalVisible}
        onClose={onModalClose}
        setRoomList={setRoomList}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginTop: StatusBar.currentHeight,
    flex: 1,
    paddingTop: 20,
    paddingBottom: 20,
  },
  row: {
    backgroundColor: "#c1c4c9",
    paddingVertical: 20,
    paddingHorizontal: 10,
    borderRadius: 15,
    minHeight: 70,
    marginVertical: 3,
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    width: "90%",
  },
  deleteIcon: {
    color: "#FF0000",
  },
  deleteBtn: {
    padding: 5,
    borderRadius: 100000,
  },
  confirmBtn: {
    borderRadius: 8,
    padding: 10,
    alignItems: "center",
    justifyContent: "center",
    width: "100%",
    marginTop: 12,
  },
  buttonsContainer: {
    alignItems: "center",
    width: "100%",
    gap: 5,
  },
  flatList: {
    alignItems: "center",
    width: "100%",
  },
});
