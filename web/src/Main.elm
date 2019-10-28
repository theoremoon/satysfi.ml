module Main exposing (main)

import Base64
import Browser
import Browser.Navigation
import Html exposing (..)
import Html.Attributes exposing (attribute, class, src, type_, value)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode as D
import Json.Encode as E
import Url exposing (Url)
import Url.Builder



-- MAIN


main =
    Browser.application
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        , onUrlRequest = OnUrlRequest
        , onUrlChange = OnUrlChange
        }



-- MODEL


type FileTree
    = File String String
    | Directory String String (List FileTree) (List FileTree)


type alias Source =
    { name : String
    , path : String
    , content : String
    }


type alias Model =
    { source : Maybe Source
    , product : Maybe String
    , fileTree : Maybe FileTree
    , sidebarFlag : Bool
    }


type Msg
    = OnUrlChange Url
    | OnUrlRequest Browser.UrlRequest
    | OpenSidebar
    | CloseSidebar
    | FileTreeRequest
    | FileTreeResult (Result Http.Error FileTree)
    | GetSourceRequest String
    | GetSourceResult (Result Http.Error Source)
    | CompileRequest
    | CompileResult (Result Http.Error String)
    | SourceUpdate String



-- DECODER


fileTreeFileDecoder : D.Decoder FileTree
fileTreeFileDecoder =
    D.map2 File
        (D.field "name" D.string)
        (D.field "path" D.string)


fileTreeDecoder : D.Decoder FileTree
fileTreeDecoder =
    D.map4 Directory
        (D.field "name" D.string)
        (D.field "path" D.string)
        (D.field "childdirs" (D.list (D.lazy (\_ -> fileTreeDecoder))))
        (D.field "children" (D.list fileTreeFileDecoder))


sourceDecoder : D.Decoder Source
sourceDecoder =
    D.map3 Source
        (D.field "name" D.string)
        (D.field "path" D.string)
        (D.field "content" D.string)



-- INIT


init : () -> Url -> Browser.Navigation.Key -> ( Model, Cmd Msg )
init _ url navKey =
    ( Model Nothing Nothing Nothing False
    , fileTreeRequest
    )



-- UPDATE


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        OpenSidebar ->
            ( { model | sidebarFlag = True }, Cmd.none )

        CloseSidebar ->
            ( { model | sidebarFlag = False }, Cmd.none )

        SourceUpdate newContent ->
            let
                newSource =
                    case model.source of
                        Just source ->
                            Just { source | content = newContent }

                        Nothing ->
                            Nothing
            in
            ( { model | source = newSource }, Cmd.none )

        FileTreeRequest ->
            ( model, fileTreeRequest )

        FileTreeResult (Result.Ok fileTreeResult) ->
            ( { model | fileTree = Just fileTreeResult }, Cmd.none )

        FileTreeResult (Result.Err _) ->
            ( model, Cmd.none )

        GetSourceRequest path ->
            ( model, getFileRequest path )

        GetSourceResult (Result.Ok newSource) ->
            ( { model | source = Just newSource }, Cmd.none )

        GetSourceResult (Result.Err _) ->
            ( model, Cmd.none )

        CompileRequest ->
            case model.source of
                Just source ->
                    ( model, compileRequest source )

                Nothing ->
                    ( model, Cmd.none )

        CompileResult (Result.Ok result) ->
            ( { model | product = Just result }, Cmd.none )

        CompileResult (Result.Err _) ->
            ( model, Cmd.none )

        _ ->
            ( model, Cmd.none )



-- CMDS


fileTreeRequest : Cmd Msg
fileTreeRequest =
    Http.get
        { url = "/filetree"
        , expect = Http.expectJson FileTreeResult fileTreeDecoder
        }


getFileRequest : String -> Cmd Msg
getFileRequest path =
    Http.get
        { url = Url.Builder.absolute [ "getfile" ] [ Url.Builder.string "filename" path ]
        , expect = Http.expectJson GetSourceResult sourceDecoder
        }


compileRequest : Source -> Cmd Msg
compileRequest source =
    Http.post
        { url = "/compile"
        , body = Http.jsonBody <| E.object [ ( "path", E.string source.path ) ]
        , expect = Http.expectString CompileResult
        }



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Browser.Document Msg
view model =
    { title = "Hello World"
    , body =
        [ div
            [ class "container"
            ]
            [ if model.sidebarFlag then
                div [ class "sidebar" ]
                    [ sidebarBackground
                        [ class "sidebarBackground"
                        , onClick CloseSidebar
                        ]
                        model
                    , fileTree
                        [ class "fileTree"
                        , onClick CloseSidebar
                        ]
                        model
                    ]

              else
                div [] []
            , div [ class "menu" ]
                [ button [ onClick OpenSidebar ] [ text "FILES" ]
                , button [ onClick CompileRequest ] [ text "COMPILE" ]
                ]
            , div [ class "main" ]
                [ satysfiEditor [ class "editor" ] model
                , productViewer [ class "viewer" ] model
                ]
            ]
        ]
    }


satysfiEditor : List (Attribute Msg) -> Model -> Html Msg
satysfiEditor attrs model =
    case model.source of
        Just source ->
            textarea
                (attrs
                    ++ [ value source.content
                       , onInput SourceUpdate
                       ]
                )
                []

        Nothing ->
            div attrs [ text "Please Load File" ]


productViewer : List (Attribute Msg) -> Model -> Html Msg
productViewer attrs model =
    case model.product of
        Just pdf ->
            embed
                (attrs
                    ++ [ src ("data:application/pdf;base64," ++ pdf)
                       ]
                )
                []

        Nothing ->
            div attrs [ text "Nothing there" ]


fileTreeImpl : FileTree -> Html Msg
fileTreeImpl filetree =
    case filetree of
        File name path ->
            li []
                [ a [ onClick (GetSourceRequest path) ] [ text name ]
                ]

        Directory name _ dirs children ->
            li []
                [ text name
                , ul [] (List.map fileTreeImpl dirs)
                , ul [] (List.map fileTreeImpl children)
                ]


fileTree : List (Attribute Msg) -> Model -> Html Msg
fileTree attrs model =
    case model.fileTree of
        Just tree ->
            ul attrs [ fileTreeImpl tree ]

        Nothing ->
            ul attrs []


sidebarBackground : List (Attribute Msg) -> Model -> Html Msg
sidebarBackground attrs model =
    div attrs []
