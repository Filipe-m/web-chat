import { PropsWithChildren, useState } from "react";
import {
  Modal,
  View,
  Text,
  Button,
  StyleSheet,
  Pressable,
  TextInput,
  KeyboardAvoidingView,
  Platform,
} from "react-native";
import DefaultBtn from "./defaultButton";
import api from "@/scripts/axios";

type Props = PropsWithChildren<{
  isVisible: boolean;
  onClose: () => void;
  setRoomList: (updater: (prev: ChatRoom[]) => ChatRoom[]) => void;
}>;

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

export default function CreateChat({ isVisible, onClose, setRoomList }: Props) {
  const [name, setName] = useState("");

  const createRoom = async () => {
    try {
      const response = await api.post<ChatRoom>("/chat", {
        name: name,
      });
      setRoomList((prev) => [...prev, response.data]);
      setName("");
    } catch (error) {
      console.log(error);
    }
    onClose();
  };

  return (
    <Modal
      animationType="slide"
      transparent={true}
      visible={isVisible}
      onRequestClose={onClose}
    >
      <Pressable style={styles.outside} onPress={onClose} />

      <View style={styles.container}>
        <View>
          <Text>Nome:</Text>

          <TextInput
            style={styles.nameInput}
            value={name}
            onChangeText={setName}
          />
        </View>
        <View>
          <View style={styles.buttonsContainer}>
            <DefaultBtn
              text="Confirmar"
              method={createRoom}
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
              text="Cancelar"
              method={onClose}
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
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  container: {
    height: "40%",
    width: "100%",
    backgroundColor: "#FFFFFF",
    borderTopRightRadius: 50,
    borderTopLeftRadius: 50,
    position: "absolute",
    bottom: 0,
    flex: 1,
    flexDirection: "column",
    justifyContent: "space-around",
    padding: 30,
    paddingTop: 20,
  },
  outside: {
    flex: 1,
  },
  nameInput: {
    borderWidth: 1,
    borderRadius: 10,
    paddingLeft: 10,
  },
  buttonsContainer: {
    marginTop: 10,
    marginBottom: 20,
    alignItems: "center",
    width: "100%",
    gap: 10,
  },
});
